package app

import (
	"log"
	"os"

	"github.com/ayo-ajayi/edutech/internal/auth"
	"github.com/ayo-ajayi/edutech/internal/db"
	"github.com/ayo-ajayi/edutech/internal/student"
	"github.com/ayo-ajayi/edutech/internal/subject"
	"github.com/ayo-ajayi/edutech/internal/tutor"
	"github.com/ayo-ajayi/edutech/internal/user"
	"github.com/ayo-ajayi/edutech/internal/utils"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	mongoDbUri := os.Getenv("MONGODB_URI")
	mongoDbName := os.Getenv("MONGODB_NAME")
	emailApiKey := os.Getenv("EMAIL_API_KEY")
	emailSenderName := os.Getenv("EMAIL_SENDER_NAME")
	emailSenderAddress := os.Getenv("EMAIL_SENDER_ADDRESS")
	verifyEmailBaseUrl := os.Getenv("BASE_URL") + "/api/v1"
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	client, err := db.MongoClient(mongoDbUri)
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx, cancel := db.DBReqContext(10)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("mongodb connected")

	verificationTokenCollection := db.NewMongoCollection(client, mongoDbName, "verification_tokens")
	verificationTokenDatabase := db.NewDatabase(verificationTokenCollection)
	verificationTokenManager := utils.NewVerificationTokenManager(verificationTokenDatabase, 60*60*24*7, 60*60*24*7)

	accessTokenCollection := db.NewMongoCollection(client, mongoDbName, "access_tokens")
	accessTokenDatabase := db.NewDatabase(accessTokenCollection)
	accessTokenManager := utils.NewTokenAccessManager(accessTokenSecret, 60*60*24*7, accessTokenDatabase)

	emailManager := utils.NewEmailManager(emailSenderAddress, emailSenderName, emailApiKey)

	studentSubjectTutorRepo := subject.NewStudentSubjectTutorRepo(db.NewDatabase(db.NewMongoCollection(client, mongoDbName, "student_subject_tutor")))
	subjectRepo := subject.NewSubjectRepo(db.NewDatabase(db.NewMongoCollection(client, mongoDbName, "subjects")))
	subjectService, err := subject.NewSubjectService(subjectRepo, "English")
	if err != nil {
		log.Fatalln("error: subject service init error: ", err.Error())
	}
	subjectController := subject.NewSubjectController(subjectService)

	tutorRepo := tutor.NewTutorRepo(db.NewDatabase(db.NewMongoCollection(client, mongoDbName, "tutors")))
	tutorService := tutor.NewTutorService(tutorRepo, verificationTokenManager, accessTokenManager, emailManager, verifyEmailBaseUrl)
	tutorController := tutor.NewTutorController(tutorService)

	studentRepo := student.NewStudentRepo(db.NewDatabase(db.NewMongoCollection(client, mongoDbName, "students")))
	studentService := student.NewStudentService(studentRepo, verificationTokenManager, accessTokenManager, emailManager, subjectRepo, tutorRepo, studentSubjectTutorRepo, verifyEmailBaseUrl)
	studentController := student.NewStudentController(studentService)

	authService := auth.NewAuthService(tutorRepo, studentRepo, subjectRepo, accessTokenManager, verificationTokenManager, emailManager, verifyEmailBaseUrl)
	authController := auth.NewAuthController(authService)

	middleware := auth.NewAuthMiddleWare(accessTokenSecret, tutorRepo, studentRepo, accessTokenManager)

	r := gin.Default()
	r.Use(jsonMiddleware(), auth.NewCors())
	r.GET("/favicon.ico", func(ctx *gin.Context) { ctx.File("./favicon.ico") })
	r.NoRoute(func(ctx *gin.Context) { ctx.JSON(404, gin.H{"error": "endpoint not found"}) })
	r.GET("/healthz", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "ok"}) })
	r.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "welcome to edutech"}) })

	api := r.Group("/api/v1")
	api.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "welcome to edutech API"}) })
	api.POST("/login", authController.Login)
	api.POST("/forgot-password", authController.ForgotPassword)
	api.POST("/reset-password", authController.ResetPassword)
	api.GET("/verify/:token", authController.Verify)
	api.DELETE("/logout", authController.Logout)

	studentRouter := api.Group("/students")
	studentRouter.POST("", studentController.SignUp)
	studentRouter.Use(middleware.Authentication(), middleware.Authorization(user.Student))
	studentRouter.GET("/profile", studentController.Profile)
	studentRouter.GET("/subjects", studentController.GetRegisteredSubjects)
	studentRouter.POST("/subjects", studentController.RegisterSubject)
	studentRouter.POST("/tutors/register", studentController.RegisterTutor)
	studentRouter.GET("/tutors", studentController.GetRegisteredTutors)

	tutorRouter := api.Group("/tutors")
	tutorRouter.POST("", tutorController.SignUp)
	tutorRouter.Use(middleware.Authentication(), middleware.Authorization(user.Tutor))
	tutorRouter.GET("/profile", tutorController.Profile)

	api.POST("/subjects", subjectController.CreateSubject)

	return r
}

func jsonMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}

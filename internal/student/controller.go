package student

import (
	"net/http"

	"github.com/ayo-ajayi/edutech/internal/user"
	"github.com/ayo-ajayi/edutech/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentController struct {
	studentService IStudentService
}

func NewStudentController(studentService IStudentService) *StudentController {
	return &StudentController{
		studentService: studentService,
	}
}

func (sc *StudentController) SignUp(c *gin.Context) {
	req := utils.SignUpReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	student := &Student{
		Id: primitive.NewObjectID(),
		User: &user.User{
			Email:     req.Email,
			Password:  req.Password,
			Firstname: req.FirstName,
			Lastname:  req.LastName,
		}}
	err := sc.studentService.SignUpStudent(student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(student, "student created successfully...check email for verification link"))
}

func (sc *StudentController) Profile(c *gin.Context) {
	id := c.MustGet("user_id").(primitive.ObjectID)
	student, err := sc.studentService.GetStudent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(student, "student retrieved successfully"))
}

//admin
// func (sc *StudentController) GetStudent(c *gin.Context) {
// 	id := c.Param("id")
// 	if id == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "id is required"}})
// 		return
// 	}

// 	studentId, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid id"}})
// 		return
// 	}
// 	student, err := sc.studentService.GetStudent(studentId)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
// 		return
// 	}
// 	c.JSON(http.StatusOK, utils.NewSuccessResponse(student, "student retrieved successfully"))
// }

func (sc *StudentController) RegisterSubject(c *gin.Context) {
	req := struct {
		SubjectId string `json:"subject_id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	subjectId, err := primitive.ObjectIDFromHex(req.SubjectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid subject id"}})
		return
	}
	userId := c.MustGet("user_id").(primitive.ObjectID)
	err = sc.studentService.RegisterSubject(subjectId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(nil, "subject registered successfully"))

}

func (sc *StudentController) GetRegisteredSubjects(c *gin.Context) {
	id := c.MustGet("user_id").(primitive.ObjectID)
	subjects, err := sc.studentService.GetRegisteredSubjects(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(subjects, "student's subjects retrieved successfully"))
}

func (sc *StudentController) RegisterTutor(c *gin.Context) {
	req := struct {
		TutorId string `json:"tutor_id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	tutorId, err := primitive.ObjectIDFromHex(req.TutorId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid tutor id"}})
		return
	}
	userId := c.MustGet("user_id").(primitive.ObjectID)
	err = sc.studentService.RegisterTutor(tutorId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(nil, "tutor registered successfully"))
}

func (sc *StudentController) GetRegisteredTutors(c *gin.Context) {
	id := c.MustGet("user_id").(primitive.ObjectID)
	tutors, err := sc.studentService.GetRegisteredTutors(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(tutors, "student's tutors retrieved successfully"))
}

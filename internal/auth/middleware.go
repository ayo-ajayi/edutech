package auth

import (
	"errors"
	"time"

	"net/http"
	"strings"

	"github.com/ayo-ajayi/edutech/internal/student"
	"github.com/ayo-ajayi/edutech/internal/tutor"
	"github.com/ayo-ajayi/edutech/internal/user"
	"github.com/ayo-ajayi/edutech/internal/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthMiddleware struct {
	accessTokenSecret  string
	tutorRepo          tutor.IMiddlewareTutorRepo
	studentRepo        student.IMiddlewareStudentRepo
	accessTokenManager utils.IMiddlewareAccessTokenManager
}

func NewAuthMiddleWare(accessTokenSecret string, tutorRepo tutor.IMiddlewareTutorRepo, studentRepo student.IMiddlewareStudentRepo, accessTokenManager utils.IMiddlewareAccessTokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		accessTokenSecret:  accessTokenSecret,
		tutorRepo:          tutorRepo,
		studentRepo:        studentRepo,
		accessTokenManager: accessTokenManager,
	}
}

func (amw *AuthMiddleware) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := amw.extractToken(c.Request)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "unauthorized: token is required"}})
			c.Abort()
			return
		}
		jwtToken, err := amw.accessTokenManager.ValidateAccessToken(token, amw.accessTokenSecret)
		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "unauthorized: invalid token"}})
				c.Abort()
				return
			}
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "unauthorized: token expired"}})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": "internal server error"}})
			c.Abort()
			return
		}
		td, err := amw.accessTokenManager.ExtractAccessTokenMetadata(jwtToken)
		if err != nil {

			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "unauthorized: invalid token"}})
			c.Abort()
			return
		}
		accessDetails, err := amw.accessTokenManager.FindAccessToken(td.AccessUuid)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "unauthorized: invalid token"}})
			c.Abort()
			return
		}
		if td.UserId != accessDetails.UserId {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "unauthorized: invalid token"}})
			c.Abort()
			return
		}
		c.Set("access_uuid", td.AccessUuid)
		c.Set("user_id", accessDetails.UserId)
		c.Next()
	}

}

func (amw *AuthMiddleware) extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	ttoken := strings.Split(token, " ")
	if len(ttoken) != 2 {
		return ""
	}
	return ttoken[1]
}

func (amw *AuthMiddleware) Authorization(role user.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.MustGet("user_id").(primitive.ObjectID)
		var currentUserRole user.Role

		if role == user.Tutor {
			tutor, err := amw.tutorRepo.GetTutor(bson.M{
				"_id": userId})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{"message": err.Error() + ": you are not authorized to access this resource"}})
				return
			}
			currentUserRole = tutor.Role
		} else if role == user.Student {
			student, err := amw.studentRepo.GetStudent(bson.M{
				"_id": userId})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{"message": err.Error() + ": you are not authorized to access this resource"}})
				return
			}
			currentUserRole = student.Role
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{"message": "invalid role"}})
			return
		}
		allowed := false
		if role == currentUserRole {
			allowed = true
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{"message": "Forbidden: You are not authorized to access this resource"}})
			return
		}
		c.Next()
	}
}

func NewCors() gin.HandlerFunc {
	cfg := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(cfg)
}

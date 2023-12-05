package tutor

import (
	"net/http"

	"github.com/ayo-ajayi/edutech/internal/user"
	"github.com/ayo-ajayi/edutech/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TutorController struct {
	tutorService ITutorService
}

func NewTutorController(tutorService ITutorService) *TutorController {
	return &TutorController{
		tutorService,
	}
}

func (tc *TutorController) SignUp(c *gin.Context) {
	req := utils.SignUpReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	tutor := &Tutor{
		Id: primitive.NewObjectID(),
		User: &user.User{
			Email:     req.Email,
			Password:  req.Password,
			Firstname: req.FirstName,
			Lastname:  req.LastName,
		}}
	err := tc.tutorService.SignUpTutor(tutor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(tutor, "tutor created successfully...check email for verification link"))
}

func (tc *TutorController) Profile(c *gin.Context) {
	id := c.MustGet("user_id").(primitive.ObjectID)
	tutor, err := tc.tutorService.GetTutor(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(tutor, "tutor retrieved successfully"))
}

//admin
// func (tc *TutorController) GetTutor(c *gin.Context) {
// 	id := c.Param("id")
// 	if id == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "id is required"}})
// 		return
// 	}

// 	tutorId, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid id"}})
// 		return
// 	}
// 	tutor, err := tc.tutorService.GetTutor(tutorId)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
// 		return
// 	}
// 	c.JSON(http.StatusOK, utils.NewSuccessResponse(tutor, "tutor retrieved successfully"))
// }

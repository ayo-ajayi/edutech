package common

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ayo-ajayi/edutech/internal/utils"
)

type UserController struct {
	userService IUserService
}

func (uc *UserController) GetTutors(c *gin.Context) {
	tutors, err := uc.userService.GetTutors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(tutors, "tutors retrieved successfully"))
}

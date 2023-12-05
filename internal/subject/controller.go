package subject

import (
	"net/http"

	"github.com/ayo-ajayi/edutech/internal/utils"
	"github.com/gin-gonic/gin"
)

type SubjectController struct {
	subjectService ISubjectService
}

func NewSubjectController(subjectService ISubjectService) *SubjectController {
	return &SubjectController{subjectService: subjectService}
}

func (sc *SubjectController) CreateSubject(c *gin.Context) {
	req := struct {
		Name string `json:"name" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	subject := &Subject{
		Name: req.Name,
	}
	if err := sc.subjectService.CreateSubject(subject); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(subject, "subject successfully created"))
}

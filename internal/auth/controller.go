package auth

import (
	"github.com/ayo-ajayi/edutech/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type AuthController struct {
	authService IAuthService
}

func NewAuthController(authService IAuthService) *AuthController {
	return &AuthController{
		authService,
	}
}
func (ac *AuthController) Verify(c *gin.Context) {
	token := c.Param("token")
	email, err := url.QueryUnescape(c.Query("email"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid email"}})
		return
	}
	if err := ac.authService.Verify(email, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(nil, "user verified successfully"))
}
func (ac *AuthController) Login(c *gin.Context) {
	req := utils.LoginReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	user, tokenDetails, err := ac.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(gin.H{
		"user":          user,
		"token_details": tokenDetails,
	}, "user successfully logged in"))
}

func (ac *AuthController) Logout(c *gin.Context) {
	accessUuid := c.GetString("access_uuid")
	err := ac.authService.Logout(accessUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(nil, "user successfully logged out"))
}

func (ac *AuthController) ForgotPassword(c *gin.Context) {
	req := struct {
		Email string `json:"email" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	err := ac.authService.ForgotPassword(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(nil, "password reset token sent successfully"))
}

func (ac *AuthController) ResetPassword(c *gin.Context) {
	req := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Token    string `json:"token" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	err := ac.authService.ResetPassword(req.Email, req.Password, req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(nil, "password reset successfully"))
}

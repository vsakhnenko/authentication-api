package routes

import (
	"authentication/internal/entities"
	"authentication/internal/middlewares"
	"authentication/internal/models"
	"authentication/internal/services"
	"authentication/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RegisterAuthRoutes(r *gin.Engine, authService services.AuthService) {
	api := r.Group("api/v1")
	{
		api.POST("/register", func(c *gin.Context) {
			register(c, authService)
		})
		api.POST("/login", func(c *gin.Context) {
			login(c, authService)
		})
		api.POST("/refresh", func(c *gin.Context) {
			refreshToken(c)
		})

		api.POST("/forgot-password", func(c *gin.Context) {
			forgotPassword(c, authService)
		})

		api.POST("/reset-password", func(c *gin.Context) {
			resetPassword(c, authService)
		})
	}

	protected := r.Group("api/v1/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", func(c *gin.Context) {
		currentUser(c, authService)
	})
}

func resetPassword(c *gin.Context, authService services.AuthService) {
	var input models.ResetPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := authService.ResetPassword(c.Request.Context(), input); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func forgotPassword(c *gin.Context, authService services.AuthService) {
	var input models.ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := authService.SendResetPassword(c.Request.Context(), input.Email)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func register(c *gin.Context, userService services.AuthService) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var u entities.User
	if err := copier.Copy(&u, &input); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := userService.Create(&u); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration success"})
}

func login(c *gin.Context, userService services.AuthService) {
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := userService.LoginCheck(input)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expiresIn, err := utils.GetTokenExpiry(accessToken)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokenResponse := models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}
	c.JSON(http.StatusOK, tokenResponse)
}

func refreshToken(c *gin.Context) {
	newAccessToken, refreshToken, err := utils.RefreshAccessToken(c)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	expiresIn, err := utils.GetTokenExpiry(newAccessToken)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenResponse := models.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}
	c.JSON(http.StatusOK, tokenResponse)
}

func currentUser(c *gin.Context, userService services.AuthService) {

	userId, err := utils.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := userService.GetByID(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

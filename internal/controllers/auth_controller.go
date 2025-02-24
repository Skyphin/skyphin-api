package controllers

import (
	"net/http"

	"skyphin-api/internal/models"
	"skyphin-api/internal/repositories"
	"skyphin-api/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService  *services.AuthService
	userService  *services.UserService
	emailService *services.EmailService
	userRepo     *repositories.UserRepository
}

func NewAuthController(authService *services.AuthService, userService *services.UserService) *AuthController {
	return &AuthController{authService: authService, userService: userService}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.userService.CreateUser(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userRepo.FindByEmail(req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	token, err := c.authService.GenerateVerificationToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	c.emailService.SendVerificationEmail(user.Email, token)

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered. Verification code sent to your email."})
}

func (c *AuthController) Verify(ctx *gin.Context) {
	var req models.VerifyAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.authService.VerifyAccount(req.Token); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Account verified successfully"})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Account not verified"})
		return
	}

	accessToken, refreshToken, err := c.authService.GenerateTokens(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

func (c *AuthController) Refresh(ctx *gin.Context) {
	var req models.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := c.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func (c *AuthController) ResetPasswordRequest(ctx *gin.Context) {
	var req models.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.authService.GenerateResetToken(req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.emailService.SendVerificationEmail(req.Email, token)

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req models.NewPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.authService.ResetPassword(&req); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"skyphin-api/internal/config"
	"skyphin-api/internal/models"
	"skyphin-api/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	authRepo *repositories.AuthRepository
	cfg      config.Config
}

func NewAuthService(userRepo *repositories.UserRepository, authRepo *repositories.AuthRepository, cfg config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, authRepo: authRepo, cfg: cfg}
}

func (s *AuthService) GenerateVerificationToken(userID string) (string, error) {
	token, err := generateRandomToken(32)
	if err != nil {
		return "", err
	}

	verificationToken := &models.VerificationToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour), // Token expires in 1 hour
	}

	if err := s.authRepo.CreateVerificationToken(verificationToken); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) VerifyAccount(token string) error {
	verificationToken, err := s.authRepo.FindVerificationToken(token)
	if err != nil || verificationToken.ExpiresAt.Before(time.Now()) {
		return errors.New("invalid or expired verification code")
	}

	user, err := s.userRepo.FindByID(verificationToken.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	user.Verified = true

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	if err := s.authRepo.DeleteVerificationToken(verificationToken.Token); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *AuthService) GenerateResetToken(email string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	token, err := generateRandomToken(32)
	if err != nil {
		return "", err
	}

	resetToken := &models.ResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour), // Token expires in 1 hour
	}

	if err := s.authRepo.CreateResetToken(resetToken); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ResetPassword(req *models.NewPasswordRequest) error {
	if req.Token == "" {
		return errors.New("empty token")
	}

	resetToken, err := s.authRepo.FindResetToken(req.Token)
	if err != nil || resetToken.ExpiresAt.Before(time.Now()) {
		return errors.New("invalid or expired token")
	}

	user, err := s.userRepo.FindByID(resetToken.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.EncryptedPassword = string(hashedPassword)

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	if err := s.authRepo.DeleteResetToken(resetToken.Token); err != nil { // Delete the token after use.
		return err
	}

	return nil
}

func (s *AuthService) GenerateTokens(user *models.User) (string, string, error) {
	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * time.Duration(s.cfg.Auth.AccessTokenExpiryMinutes)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.Auth.AccessTokenSecret))
	if err != nil {
		return "", err
	}

	accessToken := &models.AccessToken{
		UserID:    userID,
		Token:     signedToken,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(s.cfg.Auth.AccessTokenExpiryMinutes)),
	}

	if err := s.authRepo.CreateAccessToken(accessToken); err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *AuthService) generateRefreshToken(userID string) (string, error) {
	refreshTokenStr, err := generateRandomToken(64)
	if err != nil {
		return "", err
	}

	refreshToken := &models.RefreshToken{
		UserID:    userID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(s.cfg.Auth.RefreshTokenExpiryDays)),
	}

	if err := s.authRepo.CreateRefreshToken(refreshToken); err != nil {
		return "", err
	}

	return refreshTokenStr, nil
}

func (s *AuthService) RefreshAccessToken(refreshTokenStr string) (string, error) {
	refreshToken, err := s.authRepo.FindRefreshToken(refreshTokenStr)
	if err != nil || refreshToken.ExpiresAt.Before(time.Now()) {
		return "", errors.New("invalid or expired refresh token")
	}

	user, err := s.userRepo.FindByID(refreshToken.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	newAccessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func generateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

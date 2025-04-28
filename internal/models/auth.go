package models

import "time"

type VerificationToken struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	Token     string `gorm:"token"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

type AccessToken struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	Token     string `gorm:"token"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

type RefreshToken struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	Token     string `gorm:"token"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

type ResetToken struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	Token     string `gorm:"token"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

type VerifyAccountRequest struct {
	Token string `gorm:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type NewPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

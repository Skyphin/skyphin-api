package repositories

import (
	"skyphin-api/internal/models"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateVerificationToken(token *models.VerificationToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) CreateResetToken(token *models.ResetToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) CreateAccessToken(token *models.AccessToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) FindVerificationToken(token string) (*models.VerificationToken, error) {
	var at models.VerificationToken
	if err := r.db.Where("token = ?", token).First(&at).Error; err != nil {
		return nil, err
	}
	return &at, nil
}

func (r *AuthRepository) FindResetToken(token string) (*models.ResetToken, error) {
	var at models.ResetToken
	if err := r.db.Where("token = ?", token).First(&at).Error; err != nil {
		return nil, err
	}
	return &at, nil
}

func (r *AuthRepository) FindAccessToken(token string) (*models.AccessToken, error) {
	var at models.AccessToken
	if err := r.db.Where("token = ?", token).First(&at).Error; err != nil {
		return nil, err
	}
	return &at, nil
}

func (r *AuthRepository) FindRefreshToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	if err := r.db.Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *AuthRepository) DeleteVerificationToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.VerificationToken{}).Error
}

func (r *AuthRepository) DeleteResetToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.ResetToken{}).Error
}

func (r *AuthRepository) DeleteAccessToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.AccessToken{}).Error
}

func (r *AuthRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

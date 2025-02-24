package services

import "skyphin-api/internal/config"

type EmailService struct {
	cfg config.Config
}

func NewEmailService(cfg config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

func (s *EmailService) SendVerificationEmail(email string, token string) error {
	return nil
}

func (s *EmailService) SendPasswordResetEmail(email string, token string) error {
	return nil
}

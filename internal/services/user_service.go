package services

import (
	"errors"
	"skyphin-api/internal/models"
	"skyphin-api/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.repo.FindByUsername(username)
}

func (s *UserService) CreateUser(req *models.CreateUserRequest) error {
	if _, err := s.repo.FindByUsername(req.Username); err == nil {
		return errors.New("username already exists")
	}
	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:          req.Username,
		Email:             req.Email,
		EncryptedPassword: string(hashedPassword),
	}

	return s.repo.Create(user)
}

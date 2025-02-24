package main

import (
	"log"

	"skyphin-api/internal/config"
	"skyphin-api/internal/controllers"
	"skyphin-api/internal/models"
	"skyphin-api/internal/repositories"
	"skyphin-api/internal/services"
	"skyphin-api/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.AccessToken{}, &models.RefreshToken{}, &models.VerificationToken{}, &models.ResetToken{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	router := gin.Default()

	router.POST("/users", userController.CreateUser)
	router.GET("/users/:id", userController.GetUser)

	if err := router.Run(cfg.Server.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

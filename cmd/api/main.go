package main

import (
	"log"

	"skyphin-api/internal/config"
	"skyphin-api/internal/controllers"
	"skyphin-api/internal/middleware"
	"skyphin-api/internal/models"
	"skyphin-api/internal/repositories"
	"skyphin-api/internal/services"
	"skyphin-api/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	cfg := loadConfig()
	db := connectDatabase(cfg)
	migrateDatabase(db)

	userRepo, authRepo := initializeRepositories(db)
	userService, authService := initializeServices(userRepo, authRepo, cfg)
	userController, authController := initializeControllers(userService, authService)
	authMiddleware := middleware.NewAuthMiddleware(authService, cfg)

	router := setupRouter(userController, authController, authMiddleware)

	startServer(router, cfg)
}

func loadConfig() config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}

func connectDatabase(cfg config.Config) *gorm.DB {
	db, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func migrateDatabase(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}, &models.AccessToken{}, &models.RefreshToken{}, &models.VerificationToken{}, &models.ResetToken{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func initializeRepositories(db *gorm.DB) (*repositories.UserRepository, *repositories.AuthRepository) {
	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	return userRepo, authRepo
}

func initializeServices(userRepo *repositories.UserRepository, authRepo *repositories.AuthRepository, cfg config.Config) (*services.UserService, *services.AuthService) {
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, authRepo, cfg)
	return userService, authService
}

func initializeControllers(userService *services.UserService, authService *services.AuthService) (*controllers.UserController, *controllers.AuthController) {
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService, userService)
	return userController, authController
}

func setupRouter(userController *controllers.UserController, authController *controllers.AuthController, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	router := gin.Default()

	router.POST("/users", userController.CreateUser)
	router.POST("/login", authController.Login)
	router.POST("/refresh", authController.Refresh)
	router.POST("/verify", authController.Verify)
	router.POST("/reset-password-request", authController.ResetPasswordRequest)
	router.POST("/reset-password", authController.ResetPassword)

	protected := router.Group("/v1")
	protected.Use(authMiddleware.Authenticate())
	{
		// TODO
	}
	return router
}

func startServer(router *gin.Engine, cfg config.Config) {
	if err := router.Run(cfg.Server.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

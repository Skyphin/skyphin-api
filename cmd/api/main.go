package main

import (
	"log"

	"skyphin-api/graph"
	"skyphin-api/internal/config"
	"skyphin-api/internal/controllers"
	"skyphin-api/internal/middleware"
	"skyphin-api/internal/models"
	"skyphin-api/internal/repositories"
	"skyphin-api/internal/services"
	"skyphin-api/pkg/database"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	cfg := loadConfig()
	db := connectDatabase(cfg)
	if err := migrateDatabase(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	userRepo, authRepo := initializeRepositories(db)
	userService, authService := initializeServices(userRepo, authRepo, cfg)
	userController, authController := initializeControllers(userService, authService)
	authMiddleware := middleware.NewAuthMiddleware(authService, cfg)

	router := gin.Default()

	setupRestRoutes(userController, authController, authMiddleware)

	setupGraphQL(router, authMiddleware)

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

func migrateDatabase(db *gorm.DB) error {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}
	err := db.AutoMigrate(&models.User{}, &models.AccessToken{}, &models.RefreshToken{}, &models.VerificationToken{}, &models.ResetToken{})

	return err
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

func setupRestRoutes(userController *controllers.UserController, authController *controllers.AuthController, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
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

func setupGraphQL(router *gin.Engine, authMiddleware *middleware.AuthMiddleware) {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	router.GET("/graphql", gin.WrapH(playground.Handler("GraphQL playground", "/query")))

	router.POST("/query", authMiddleware.Authenticate(), gin.WrapH(srv))
}

func startServer(router *gin.Engine, cfg config.Config) {
	if err := router.Run(cfg.Server.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

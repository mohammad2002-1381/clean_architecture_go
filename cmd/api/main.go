package main

import (
	"log"
	"os"

	"go-ca/internal/api/admincontroller"
	"go-ca/internal/api/middleware"
	"go-ca/internal/api/usercontroller"
	"go-ca/internal/app/notification"
	"go-ca/internal/app/service"
	"go-ca/internal/domain"
	"go-ca/internal/infra"
	"go-ca/internal/infra/postgres"

	_ "go-ca/cmd/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Clean Architecture Go API
// @version 1.0
// @description This is the API documentation for the Go application.
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to environment variables")
	}

	// Fetch environment variables
	dbConfig := os.Getenv("DB_CONFIG")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")

	db, err := postgres.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := infra.NewBaseRepository[*domain.User, uint](db)
	productRepo := infra.NewBaseRepository[*domain.Product, uint](db)
	tokenRepo := infra.NewBaseRepository[*domain.Token, uint](db)
	jwtService := service.NewJwtService(jwtSecret)
	passwordService := service.NewPasswordService()
	currentUserService := service.NewCurrentUserService()

	RegisterEvents()

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	authMiddleware := middleware.AuthMiddleware(jwtService)

	v1Group := router.Group("/api/v1")

	usercontroller.RegisterUserController(
		&userRepo,
		&tokenRepo,
		&productRepo,
		v1Group,
		jwtService,
		passwordService,
		currentUserService,
		authMiddleware,
	)

	admincontroller.RegisterAdminController(v1Group, &userRepo, authMiddleware)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func RegisterEvents() {
	productEventHandler := notification.NewCreateProductCommandHandler()
	domain.RegisterHandler[*domain.ProductCreatedEvent](&productEventHandler)
}

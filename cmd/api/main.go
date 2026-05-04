package main

import (
	"log"

	"go-ca/internal/api/admincontroller"
	"go-ca/internal/api/middleware"
	"go-ca/internal/api/usercontroller"
	"go-ca/internal/app/service"
	"go-ca/internal/domain"
	"go-ca/internal/infra"
	"go-ca/internal/infra/postgres"

	_ "go-ca/cmd/docs"

	"github.com/gin-gonic/gin"

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
	db, err := postgres.NewDatabase("host=87.248.131.253 user=shahiapp_user password=shahiapp_password dbname=clean_architecture_go port=5432")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := infra.NewBaseRepository[domain.User, uint](db)
	productRepo := infra.NewBaseRepository[domain.Product, uint](db)
	jwtService := service.NewJwtService("your-super-secret-key")
	passwordService := service.NewPasswordService()
	currentUserService := service.NewCurrentUserService()

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	authMiddleware := middleware.AuthMiddleware(jwtService)

	v1Group := router.Group("/api/v1")

	usercontroller.RegisterUserController(
		&userRepo,
		&productRepo,
		v1Group,
		jwtService,
		passwordService,
		currentUserService,
		authMiddleware,
	)

	admincontroller.RegisterAdminController(v1Group, &userRepo, authMiddleware)

	if err := router.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

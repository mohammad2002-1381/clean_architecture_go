package main

import (
	"log"

	"go-ca/internal/api"
	"go-ca/internal/api/middleware"
	"go-ca/internal/app/productapp"
	"go-ca/internal/app/service"
	"go-ca/internal/app/userapp"
	"go-ca/internal/domain"
	"go-ca/internal/infra"
	"go-ca/internal/infra/postgres"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

	authMiddleware := middleware.AuthMiddleware(jwtService)

	v1Group := router.Group("/api/v1")

	registerProductModule(db, v1Group, &productRepo, authMiddleware, currentUserService)
	registerUserModule(db, &userRepo, v1Group, jwtService, passwordService, currentUserService, authMiddleware)

	if err := router.Run(":9000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func registerProductModule(db *gorm.DB,
	routerGroup *gin.RouterGroup,
	productRepo domain.BaseRepository[domain.Product, uint],
	authMiddleware gin.HandlerFunc,
	currentUserService service.CurrentUserService,
) {
	createHandler := productapp.NewCreateProductCommandHandler(productRepo, currentUserService)
	updateHandler := productapp.NewUpdateProductCommandHandler(productRepo, currentUserService)
	getByIdHandler := productapp.NewGetProductByIdQueryHandler(productRepo, currentUserService)
	getListHandler := productapp.NewGetProductListQueryHandler(productRepo, currentUserService)

	productController := api.NewProductController(
		&createHandler,
		&updateHandler,
		&getByIdHandler,
		&getListHandler,
	)

	productController.RegisterRoutes(routerGroup, authMiddleware)
}

func registerUserModule(
	db *gorm.DB,
	userRepo domain.BaseRepository[domain.User, uint],
	routerGroup *gin.RouterGroup,
	jwtSvc service.JWTService,
	passSvc service.PasswordService,
	currentUserSvc service.CurrentUserService,
	authMiddleware gin.HandlerFunc,
) {
	registerHandler := userapp.NewRegisterUserCommandHandler(userRepo, passSvc, jwtSvc)
	loginHandler := userapp.NewLoginUserCommandHandler(userRepo, passSvc, jwtSvc)
	getUserHandler := userapp.NewGetUserQueryHandler(userRepo, currentUserSvc)

	// Initialize Controller
	userController := api.NewUserController(
		&registerHandler,
		&loginHandler,
		&getUserHandler,
		currentUserSvc,
	)

	userController.RegisterRoutes(routerGroup, authMiddleware)
}

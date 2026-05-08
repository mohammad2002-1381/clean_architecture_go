package usercontroller

import (
	"go-ca/internal/app/productapp"
	"go-ca/internal/app/service"
	"go-ca/internal/app/userapp"
	"go-ca/internal/domain"

	"github.com/gin-gonic/gin"
)

func RegisterUserController(
	userRepo domain.BaseRepository[*domain.User, uint],
	tokenRepo domain.BaseRepository[*domain.Token, uint],
	productRepo domain.BaseRepository[*domain.Product, uint],
	routerGroup *gin.RouterGroup,
	jwtSvc service.JWTService,
	passSvc service.PasswordService,
	currentUserSvc service.CurrentUserService,
	authMiddleware gin.HandlerFunc,
) {
	registerProductModule(routerGroup, productRepo, authMiddleware, currentUserSvc)
	registerUserModule(userRepo, tokenRepo, routerGroup, jwtSvc, passSvc, currentUserSvc, authMiddleware)
}

func registerProductModule(routerGroup *gin.RouterGroup,
	productRepo domain.BaseRepository[*domain.Product, uint],
	authMiddleware gin.HandlerFunc,
	currentUserService service.CurrentUserService,
) {
	createHandler := productapp.NewCreateProductCommandHandler(productRepo, currentUserService)
	updateHandler := productapp.NewUpdateProductCommandHandler(productRepo, currentUserService)
	getByIdHandler := productapp.NewGetProductByIdQueryHandler(productRepo, currentUserService)
	getListHandler := productapp.NewGetProductListQueryHandler(productRepo, currentUserService)

	productController := newProductController(
		&createHandler,
		&updateHandler,
		&getByIdHandler,
		&getListHandler,
	)

	productController.registerRoutes(routerGroup, authMiddleware)
}

func registerUserModule(
	userRepo domain.BaseRepository[*domain.User, uint],
	tokenRepo domain.BaseRepository[*domain.Token, uint],
	routerGroup *gin.RouterGroup,
	jwtSvc service.JWTService,
	passSvc service.PasswordService,
	currentUserSvc service.CurrentUserService,
	authMiddleware gin.HandlerFunc,
) {
	registerHandler := userapp.NewRegisterUserCommandHandler(userRepo, passSvc, jwtSvc)
	loginHandler := userapp.NewLoginUserCommandHandler(userRepo, passSvc, jwtSvc)
	getUserHandler := userapp.NewGetUserQueryHandler(userRepo, currentUserSvc)
	refreshTokenHandler := userapp.NewRefreshTokenCommandHandler(tokenRepo, currentUserSvc, jwtSvc)

	userController := newUserController(
		&registerHandler,
		&loginHandler,
		&getUserHandler,
		&refreshTokenHandler,
		currentUserSvc,
	)

	userController.registerRoutes(routerGroup, authMiddleware)
}

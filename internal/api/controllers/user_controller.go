// internal/api/controllers/user_controller.go
package controllers

import (
	"clean_architecture_go/internal/api/middlewares"
	"clean_architecture_go/internal/app"
	"clean_architecture_go/internal/app/features/users/commands"
	"clean_architecture_go/internal/app/features/users/dto"
	"clean_architecture_go/internal/app/features/users/queries"
	"clean_architecture_go/internal/app/mappers"
	"clean_architecture_go/internal/domain"
	"clean_architecture_go/internal/infra"
	"clean_architecture_go/internal/infra/repositories"
	"clean_architecture_go/internal/pkg/mediatr"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Login    mediatr.RequestHandler[commands.LoginUserCommand, dto.AuthDto]
	Register mediatr.RequestHandler[commands.RegisterUserCommand, dto.AuthDto]
	GetUser  mediatr.RequestHandler[queries.GetUserQuery, dto.UserDTO]
}

func RegisterUserRoutes(rg *gin.RouterGroup) {
	userRepo := repositories.NewUserRepository()
	mapper := mappers.NewUserMapper()
	tokenRepo := infra.NewBaseRepository[domain.Token, int32]()
	auth := middlewares.AuthMiddleware(app.GetTokenService())

	loginHandler := commands.NewLoginUserCommandHandler(mapper, userRepo, tokenRepo)
	registerHandler := commands.NewRegisterUserCommandHandler(mapper, userRepo, tokenRepo)
	getUserHandler := queries.NewGetUserQueryHandler(mapper, userRepo)

	controller := &UserController{
		Login:    &loginHandler,
		Register: &registerHandler,
		GetUser:  &getUserHandler,
	}

	public := rg.Group("/user")
	{
		public.POST("/login", controller.login)
		public.POST("/register", controller.register)
	}

	protected := rg.Group("/user")
	protected.Use(auth)
	{
		protected.GET("", controller.getUser)
	}
}

func (c *UserController) login(ctx *gin.Context) {
	var cmd commands.LoginUserCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := c.Login.Handle(ctx.Request.Context(), cmd)
	middlewares.Json(ctx, result, err)
}

func (c *UserController) register(ctx *gin.Context) {
	var cmd commands.RegisterUserCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := c.Register.Handle(ctx.Request.Context(), cmd)
	middlewares.Json(ctx, result, err)
}

func (c *UserController) getUser(ctx *gin.Context) {
	var query queries.GetUserQuery
	result, err := c.GetUser.Handle(ctx.Request.Context(), query)
	middlewares.Json(ctx, result, err)
}

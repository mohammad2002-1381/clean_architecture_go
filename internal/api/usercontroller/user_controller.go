package usercontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-ca/internal/app/service"
	"go-ca/internal/app/userapp"
)

type UserController struct {
	registerHandler     *userapp.RegisterUserCommandHandler
	loginHandler        *userapp.LoginUserCommandHandler
	getUserHandler      *userapp.GetUserQueryHandler
	refreshTokenHandler *userapp.RefreshTokenCommandHandler
	currentUserService  service.CurrentUserService
}

func newUserController(
	register *userapp.RegisterUserCommandHandler,
	login *userapp.LoginUserCommandHandler,
	getUser *userapp.GetUserQueryHandler,
	refreshTokenHandler *userapp.RefreshTokenCommandHandler,
	currentUserSvc service.CurrentUserService,
) *UserController {
	return &UserController{
		registerHandler:    register,
		loginHandler:       login,
		getUserHandler:     getUser,
		currentUserService: currentUserSvc,
	}
}

func (uc *UserController) registerRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) *gin.RouterGroup {
	users := router.Group("/users")
	{
		users.POST("/register", uc.Register)
		users.POST("/login", uc.Login)

		protected := users.Group("")
		protected.Use(authMiddleware)
		{
			protected.GET("", uc.GetCurrentUser)
			protected.POST("/refresh_token", uc.RefreshToken)
		}
	}

	return router
}

// Register godoc
// @Summary      Register a user
// @Description  Creates a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body userapp.RegisterUserCommand true "User registration info"
// @Success      201  {object}  userapp.AuthDTO
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /api/v1/users/register [post]
func (uc *UserController) Register(c *gin.Context) {
	var cmd userapp.RegisterUserCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	dto, err := uc.registerHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates a user and returns an access token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body userapp.LoginUserCommand true "Login credentials"
// @Success      200  {object}  userapp.AuthDTO
// @Failure      400  {object}  map[string]string "error"
// @Failure      401  {object}  map[string]string "error"
// @Router       /api/v1/users/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var cmd userapp.LoginUserCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	dto, err := uc.loginHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials or " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto)
}

// GetCurrentUser godoc
// @Summary      Get current user
// @Description  Retrieves the profile of the currently authenticated user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request query userapp.GetUserQuery true "Query context"
// @Success      200  {object}  userapp.UserDTO
// @Failure      400  {object}  map[string]string "error"
// @Failure      404  {object}  map[string]string "error"
// @Router       /api/v1/users [get]
func (uc *UserController) GetCurrentUser(c *gin.Context) {
	var query userapp.GetUserQuery

	// if err := c.ShouldBindJSON(&query); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
	// 	return
	// }

	dto, err := uc.getUserHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, dto)
}

// Login godoc
// @Summary      Refresh user token
// @Description  Refresh user token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body userapp.RefreshTokenCommand true "RefreshToken credentials"
// @Success      200  {object}  userapp.TokenDTO
// @Failure      400  {object}  map[string]string "error"
// @Failure      401  {object}  map[string]string "error"
// @Router       /api/v1/users/refresh_token [post]
func (uc *UserController) RefreshToken(c *gin.Context) {
	var cmd userapp.RefreshTokenCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	dto, err := uc.refreshTokenHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials or " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto)
}

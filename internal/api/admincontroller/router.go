package admincontroller

import (
	"go-ca/internal/api/middleware"
	"go-ca/internal/app/userapp"
	"go-ca/internal/domain"

	"github.com/gin-gonic/gin"
)

func RegisterAdminController(
	routerGroup *gin.RouterGroup,
	userRepo domain.BaseRepository[domain.User, uint],
	authMiddleware gin.HandlerFunc,
) {
	roleMiddleware := middleware.RequireRole("admin")
	registerUserModule(userRepo, routerGroup, authMiddleware, roleMiddleware)
}

func registerUserModule(
	userRepo domain.BaseRepository[domain.User, uint],
	router *gin.RouterGroup, authMiddleware gin.HandlerFunc, roleMiddleware gin.HandlerFunc,
) {
	getUserListHandler := userapp.NewGetUsersListQueryHandler(userRepo)
	controller := newAdminUserController(&getUserListHandler)

	controller.RegisterRoutes(router, authMiddleware, roleMiddleware)
}

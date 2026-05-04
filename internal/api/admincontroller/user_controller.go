package admincontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-ca/internal/app/userapp"
	_ "go-ca/internal/domain"
)

type AdminUserController struct {
	getUsersListHandler *userapp.GetUsersListQueryHandler
}

func newAdminUserController(
	getUsersList *userapp.GetUsersListQueryHandler,
) *AdminUserController {
	return &AdminUserController{
		getUsersListHandler: getUsersList,
	}
}

func (ac *AdminUserController) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, roleMiddleware gin.HandlerFunc) *gin.RouterGroup {
	adminUsers := router.Group("/admin/users")
	adminUsers.Use(authMiddleware, roleMiddleware)
	{
		adminUsers.GET("", ac.GetUsers)
	}

	return router
}

// GetUsers godoc
// @Summary      Get a list of users
// @Description  Retrieves a paginated list of users (Admin only)
// @Tags         admin-users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        query query userapp.GetUsersListQuery true "Filter and Pagination params"
// @Success      200  {object}  domain.PaginatedResult[userapp.UserDTO]
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /api/v1/admin/users [get]
func (ac *AdminUserController) GetUsers(c *gin.Context) {
	var query userapp.GetUsersListQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	result, err := ac.getUsersListHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

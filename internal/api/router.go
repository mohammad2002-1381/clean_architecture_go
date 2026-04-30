package api

import (
	"clean_architecture_go/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func InjectRouters(rg *gin.RouterGroup) {
	controllers.RegisterUserRoutes(rg)
	controllers.RegisterProductRoutes(rg)
}

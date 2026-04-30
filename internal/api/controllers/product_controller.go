// internal/api/controllers/product_controller.go
package controllers

import (
	"clean_architecture_go/internal/api/middlewares"
	"clean_architecture_go/internal/app"
	dto "clean_architecture_go/internal/app/features/products/dto"
	"clean_architecture_go/internal/app/features/products/queries"
	"clean_architecture_go/internal/domain"
	"clean_architecture_go/internal/domain/extensions"
	"clean_architecture_go/internal/infra"
	"clean_architecture_go/internal/pkg/mediatr"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	GetProducts mediatr.RequestHandler[queries.GetProductsQuery, extensions.PagedResult[*dto.ProductDTO]]
}

func RegisterProductRoutes(rg *gin.RouterGroup) {
	productRepo := infra.NewBaseRepository[domain.Product, int32]()
	auth := middlewares.AuthMiddleware(app.GetTokenService())

	getProductsHandler := queries.NewGetProductsQueryHandler(productRepo)

	controller := &ProductController{
		GetProducts: getProductsHandler,
	}

	protected := rg.Group("/product")
	protected.Use(auth)
	{
		protected.GET("", controller.getProducts)
	}
}

func (c *ProductController) getProducts(ctx *gin.Context) {
	var query queries.GetProductsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := c.GetProducts.Handle(ctx.Request.Context(), query)
	middlewares.Json(ctx, result, err)
}

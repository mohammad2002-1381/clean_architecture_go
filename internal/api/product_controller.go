package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-ca/internal/app/productapp"
)

type ProductController struct {
	createHandler  *productapp.CreateProductCommandHandler
	updateHandler  *productapp.UpdateProductCommandHandler
	getByIdHandler *productapp.GetProductByIdQueryHandler
	getListHandler *productapp.GetProductListQueryHandler
}

func NewProductController(
	create *productapp.CreateProductCommandHandler,
	update *productapp.UpdateProductCommandHandler,
	getById *productapp.GetProductByIdQueryHandler,
	getList *productapp.GetProductListQueryHandler,
) *ProductController {
	return &ProductController{
		createHandler:  create,
		updateHandler:  update,
		getByIdHandler: getById,
		getListHandler: getList,
	}
}

func (pc *ProductController) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) *gin.RouterGroup {
	products := router.Group("/products")

	products.Use(authMiddleware)
	{
		products.POST("", pc.Create)
		products.PUT("/:id", pc.Update)
		products.GET("/:id", pc.GetById)
		products.GET("", pc.GetList)
	}

	return router
}

// Create godoc
// @Summary      Create a product
// @Description  Adds a new product to the system
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body productapp.CreateProductCommand true "Product details"
// @Success      201  {object}  productapp.ProductDTO
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /products [post]
func (pc *ProductController) Create(c *gin.Context) {
	var cmd productapp.CreateProductCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	dto, err := pc.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto)
}

// Update godoc
// @Summary      Update a product
// @Description  Updates an existing product's information
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int true "Product ID"
// @Param        request body productapp.UpdateProductCommand true "Updated product details"
// @Success      200  {object}  productapp.ProductDTO
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /products/{id} [put]
func (pc *ProductController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	var cmd productapp.UpdateProductCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	cmd.ID = uint(id)

	dto, err := pc.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto)
}

// GetById godoc
// @Summary      Get a product by ID
// @Description  Retrieves a specific product by its ID
// @Tags         products
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int true "Product ID"
// @Success      200  {object}  productapp.ProductDTO
// @Failure      400  {object}  map[string]string "error"
// @Failure      404  {object}  map[string]string "error"
// @Router       /products/{id} [get]
func (pc *ProductController) GetById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	query := productapp.GetProductByIdQuery{ID: uint(id)}

	dto, err := pc.getByIdHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, dto)
}

// GetList godoc
// @Summary      List products
// @Description  Retrieves a paginated list of products based on query parameters
// @Tags         products
// @Produce      json
// @Security     BearerAuth
// @Param        query query productapp.GetProductListQuery true "Filter and Pagination params"
// @Success      200  {array}   productapp.ProductDTO
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /products [get]
func (pc *ProductController) GetList(c *gin.Context) {
	var query productapp.GetProductListQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	dtos, err := pc.getListHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos)
}

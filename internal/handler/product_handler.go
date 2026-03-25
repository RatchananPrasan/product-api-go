package handler

import (
	"errors"
	"net/http"

	"github.com/example/product-api/internal/domain"
	"github.com/example/product-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	usecase domain.ProductUsecase
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(uc domain.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: uc}
}

// RegisterRoutes registers all product routes
func (h *ProductHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/product", h.CreateProduct)
	r.PATCH("/product/:id", h.PatchProduct)
}

// CreateProduct godoc
// @Summary      Create a product
// @Description  Create a new product with name, price, optional description and sale_price
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CreateProductInput  true  "Product payload"
// @Success      201   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /product [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input domain.CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("INVALID_REQUEST"))
		return
	}

	product, err := h.usecase.CreateProduct(&input)
	if err != nil {
		if errors.Is(err, domain.ErrSalePriceMustBeLessThanPrice) {
			c.JSON(http.StatusBadRequest, response.Error("SALE_PRICE_MUST_BE_LESS_THAN_PRICE"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusCreated, response.Success(gin.H{
		"data1": product.ID.String(),
		"data2": product.Name,
	}))
}

// PatchProduct godoc
// @Summary      Patch a product
// @Description  Partially update a product — only provided fields are updated
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id    path      string                     true  "Product ID"
// @Param        body  body      domain.PatchProductInput   true  "Patch payload"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /product/{id} [patch]
func (h *ProductHandler) PatchProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("INVALID_ID"))
		return
	}

	var input domain.PatchProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("INVALID_REQUEST"))
		return
	}

	err = h.usecase.PatchProduct(id, &input)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, response.Error("PRODUCT_NOT_FOUND"))
			return
		}
		if errors.Is(err, domain.ErrSalePriceMustBeLessThanPrice) {
			c.JSON(http.StatusBadRequest, response.Error("SALE_PRICE_MUST_BE_LESS_THAN_PRICE"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

package handlers

import (
	"net/http"

	"ecommerce/internal/models"
	"ecommerce/internal/pkg/logger"
	"ecommerce/internal/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProductHandler struct {
	productRepo *repository.ProductRepository
	logger      *zap.Logger
}

func NewProductHandler(productRepo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{
		productRepo: productRepo,
		logger:      logger.GetLogger(),
	}
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.productRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productRepo.Create(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product details
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Product details"
// @Security Bearer
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		h.logger.Error("Failed to bind product JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID := c.Param("id")
	if err := h.productRepo.Update(productID, &product); err != nil {
		h.logger.Error("Failed to update product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	if err := h.productRepo.Delete(productID); err != nil {
		h.logger.Error("Failed to delete product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

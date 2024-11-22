package handlers

import (
	"net/http"

	"ecommerce/internal/models"
	"ecommerce/internal/pkg/logger"
	"ecommerce/internal/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WarehouseHandler struct {
	warehouseRepo *repository.WarehouseRepository
	productRepo   *repository.ProductRepository
	logger        *zap.Logger
}

func NewWarehouseHandler(warehouseRepo *repository.WarehouseRepository, productRepo *repository.ProductRepository) *WarehouseHandler {
	return &WarehouseHandler{
		warehouseRepo: warehouseRepo,
		productRepo:   productRepo,
		logger:        logger.GetLogger(),
	}
}

// TransferStock godoc
// @Summary Transfer stock between warehouses
// @Description Transfer product stock from one warehouse to another
// @Tags warehouse
// @Accept json
// @Produce json
// @Param transfer body models.StockTransfer true "Stock transfer details"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /warehouse/transfer [post]
func (h *WarehouseHandler) TransferStock(c *gin.Context) {
	var transfer models.StockTransfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		h.logger.Error("Failed to bind transfer JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.warehouseRepo.TransferStock(&transfer); err != nil {
		h.logger.Error("Failed to transfer stock",
			zap.String("from_warehouse", transfer.FromWarehouse.Hex()),
			zap.String("to_warehouse", transfer.ToWarehouse.Hex()),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Stock transferred successfully",
		zap.String("product_id", transfer.ProductID.Hex()),
		zap.Int("quantity", transfer.Quantity))
	c.JSON(http.StatusOK, gin.H{"message": "Stock transferred successfully"})
}

// UpdateStatus godoc
// @Summary Update warehouse status
// @Description Update the status of a warehouse
// @Tags warehouse
// @Accept json
// @Produce json
// @Param id path string true "Warehouse ID"
// @Param status body object true "Status update"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /warehouse/{id}/status [put]
func (h *WarehouseHandler) UpdateStatus(c *gin.Context) {
	warehouseID := c.Param("id")
	var status struct {
		Status models.WarehouseStatus `json:"status"`
	}

	if err := c.ShouldBindJSON(&status); err != nil {
		h.logger.Error("Failed to bind status JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.warehouseRepo.UpdateStatus(warehouseID, status.Status); err != nil {
		h.logger.Error("Failed to update warehouse status",
			zap.String("warehouse_id", warehouseID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Warehouse status updated successfully",
		zap.String("warehouse_id", warehouseID),
		zap.String("status", string(status.Status)))
	c.JSON(http.StatusOK, gin.H{"message": "Warehouse status updated successfully"})
}

// ListWarehouses godoc
// @Summary List all warehouses
// @Description Get all warehouses
// @Tags warehouses
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Warehouse
// @Router /warehouse [get]
func (h *WarehouseHandler) ListWarehouses(c *gin.Context) {
	warehouses, err := h.warehouseRepo.List()
	if err != nil {
		h.logger.Error("Failed to fetch warehouses", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, warehouses)
}

// GetProductStock godoc
// @Summary Get product stock in warehouse
// @Description Get product stock level in specific warehouse
// @Tags warehouse
// @Accept json
// @Produce json
// @Param warehouseId path string true "Warehouse ID"
// @Param productId path string true "Product ID"
// @Security Bearer
// @Success 200 {object} map[string]int
// @Router /warehouse/{warehouseId}/products/{productId} [get]
func (h *WarehouseHandler) GetProductStock(c *gin.Context) {
	warehouseID := c.Param("warehouseId")
	productID := c.Param("productId")

	stock, err := h.productRepo.GetStockInWarehouse(productID, warehouseID)
	if err != nil {
		h.logger.Error("Failed to get product stock",
			zap.String("warehouse_id", warehouseID),
			zap.String("product_id", productID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stock": stock})
}

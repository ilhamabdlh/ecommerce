package handler

import (
	"ecommerce/warehouse-service/internal/domain"
	"ecommerce/warehouse-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	warehouseUsecase domain.WarehouseUsecase
}

func NewWarehouseHandler(usecase domain.WarehouseUsecase) *WarehouseHandler {
	return &WarehouseHandler{
		warehouseUsecase: usecase,
	}
}

// @Summary Create new warehouse
// @Description Create a new warehouse
// @Tags warehouses
// @Accept json
// @Produce json
// @Param warehouse body models.Warehouse true "Warehouse Info"
// @Success 201 {object} models.Warehouse
// @Router /warehouses [post]
func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var warehouse models.Warehouse
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.warehouseUsecase.CreateWarehouse(c.Request.Context(), &warehouse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, warehouse)
}

// @Summary Get warehouse by ID
// @Description Get warehouse details by ID
// @Tags warehouses
// @Produce json
// @Param id path string true "Warehouse ID"
// @Success 200 {object} models.Warehouse
// @Router /warehouses/{id} [get]
func (h *WarehouseHandler) GetWarehouse(c *gin.Context) {
	id := c.Param("id")
	warehouse, err := h.warehouseUsecase.GetWarehouse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, warehouse)
}

// @Summary Update stock
// @Description Update product stock in warehouse
// @Tags warehouses
// @Accept json
// @Produce json
// @Param id path string true "Warehouse ID"
// @Param update body map[string]interface{} true "Stock Update"
// @Success 200 {string} string "Stock updated successfully"
// @Router /warehouses/{id}/stock [put]
func (h *WarehouseHandler) UpdateStock(c *gin.Context) {
	warehouseID := c.Param("id")
	var request struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.warehouseUsecase.UpdateStock(c.Request.Context(), warehouseID, request.ProductID, request.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}

// @Summary Transfer stock
// @Description Transfer stock between warehouses
// @Tags warehouses
// @Accept json
// @Produce json
// @Param transfer body models.StockTransfer true "Transfer Details"
// @Success 200 {string} string "Stock transferred successfully"
// @Router /warehouses/transfer [post]
func (h *WarehouseHandler) TransferStock(c *gin.Context) {
	var transfer models.StockTransfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.warehouseUsecase.TransferStock(c.Request.Context(), &transfer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock transferred successfully"})
}

// @Summary Update warehouse status
// @Description Activate or deactivate warehouse
// @Tags warehouses
// @Accept json
// @Produce json
// @Param id path string true "Warehouse ID"
// @Param status path string true "Status (activate/deactivate)"
// @Success 200 {string} string "Status updated successfully"
// @Router /warehouses/{id}/{status} [put]
func (h *WarehouseHandler) UpdateWarehouseStatus(c *gin.Context) {
	id := c.Param("id")
	status := c.Param("status")

	var err error
	switch status {
	case "activate":
		err = h.warehouseUsecase.ActivateWarehouse(c.Request.Context(), id)
	case "deactivate":
		err = h.warehouseUsecase.DeactivateWarehouse(c.Request.Context(), id)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Warehouse status updated successfully"})
}

// @Summary Get all warehouses
// @Description Get list of all warehouses
// @Tags warehouses
// @Produce json
// @Success 200 {array} models.Warehouse
// @Router /warehouses [get]
func (h *WarehouseHandler) GetAllWarehouses(c *gin.Context) {
	warehouses, err := h.warehouseUsecase.GetAllWarehouses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, warehouses)
}

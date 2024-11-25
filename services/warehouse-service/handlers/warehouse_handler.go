package handlers

import (
	"ecommerce/warehouse-service/models"
	"ecommerce/warehouse-service/repository"
	"ecommerce/warehouse-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	warehouseRepo repository.WarehouseRepository
}

func NewWarehouseHandler(repo repository.WarehouseRepository) *WarehouseHandler {
	return &WarehouseHandler{
		warehouseRepo: repo,
	}
}

// @Summary Create new warehouse
// @Description Create a new warehouse
// @Tags warehouses
// @Accept json
// @Produce json
// @Param warehouse body models.Warehouse true "Warehouse Info"
// @Success 201 {object} models.Warehouse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/warehouses [post]
func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var warehouse models.Warehouse
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		utils.Logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(warehouse); err != nil {
		utils.Logger.Errorf("Validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.warehouseRepo.Create(c.Request.Context(), &warehouse); err != nil {
		utils.Logger.Errorf("Failed to create warehouse: %v", err)
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
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/warehouses/{id} [get]
func (h *WarehouseHandler) GetWarehouse(c *gin.Context) {
	id := c.Param("id")
	warehouse, err := h.warehouseRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.Logger.Errorf("Failed to get warehouse: %v", err)
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
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/warehouses/{id}/stock [put]
func (h *WarehouseHandler) UpdateStock(c *gin.Context) {
	warehouseID := c.Param("id")
	var request struct {
		ProductID string `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.warehouseRepo.UpdateStock(c.Request.Context(), warehouseID, request.ProductID, request.Quantity)
	if err != nil {
		utils.Logger.Errorf("Failed to update stock: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}

// @Summary Transfer stock between warehouses
// @Description Transfer stock from one warehouse to another
// @Tags warehouses
// @Accept json
// @Produce json
// @Param transfer body models.StockTransfer true "Transfer Details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/warehouses/transfer [post]
func (h *WarehouseHandler) TransferStock(c *gin.Context) {
	var transfer models.StockTransfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		utils.Logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.warehouseRepo.TransferStock(c.Request.Context(), &transfer); err != nil {
		utils.Logger.Errorf("Failed to transfer stock: %v", err)
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
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/warehouses/{id}/{status} [put]
func (h *WarehouseHandler) UpdateWarehouseStatus(c *gin.Context) {
	id := c.Param("id")
	status := c.Param("status")

	if status != "activate" && status != "deactivate" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	newStatus := "active"
	if status == "deactivate" {
		newStatus = "inactive"
	}

	if err := h.warehouseRepo.UpdateStatus(c.Request.Context(), id, newStatus); err != nil {
		utils.Logger.Errorf("Failed to update warehouse status: %v", err)
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
// @Failure 500 {object} map[string]string
// @Router /api/v1/warehouses [get]
func (h *WarehouseHandler) GetAllWarehouses(c *gin.Context) {
	warehouses, err := h.warehouseRepo.GetAll(c.Request.Context())
	if err != nil {
		utils.Logger.Errorf("Failed to get warehouses: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, warehouses)
}

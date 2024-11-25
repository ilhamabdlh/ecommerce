package handlers

import (
	"net/http"

	"ecommerce/shop-service/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopHandler struct {
	shopService models.ShopService
}

func NewShopHandler(shopService models.ShopService) *ShopHandler {
	return &ShopHandler{
		shopService: shopService,
	}
}

// CreateShop godoc
// @Summary Create a new shop
// @Description Create a new shop with the input payload
// @Tags shops
// @Accept  json
// @Produce  json
// @Param shop body models.Shop true "Create shop"
// @Success 201 {object} models.Shop
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /shops [post]
func (h *ShopHandler) CreateShop(c *gin.Context) {
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.shopService.CreateShop(&shop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shop)
}

// GetShop godoc
// @Summary Get a shop by ID
// @Description Get details of a specific shop by its ID
// @Tags shops
// @Accept  json
// @Produce  json
// @Param id path string true "Shop ID"
// @Success 200 {object} models.Shop
// @Failure 400 {object} map[string]interface{} "Invalid ID format"
// @Failure 404 {object} map[string]interface{} "Shop not found"
// @Router /shops/{id} [get]
func (h *ShopHandler) GetShop(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	shop, err := h.shopService.GetShop(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
		return
	}

	c.JSON(http.StatusOK, shop)
}

// ListShops godoc
// @Summary List all shops
// @Description Get a list of all shops
// @Tags shops
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Shop
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /shops [get]
func (h *ShopHandler) ListShops(c *gin.Context) {
	shops, err := h.shopService.ListShops()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shops)
}

// UpdateShop godoc
// @Summary Update a shop
// @Description Update a shop's details by its ID
// @Tags shops
// @Accept  json
// @Produce  json
// @Param id path string true "Shop ID"
// @Param shop body models.Shop true "Update shop"
// @Success 200 {object} models.Shop
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /shops/{id} [put]
func (h *ShopHandler) UpdateShop(c *gin.Context) {
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	shop.ID = id

	if err := h.shopService.UpdateShop(&shop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shop)
}

// DeleteShop godoc
// @Summary Delete a shop
// @Description Delete a shop by its ID
// @Tags shops
// @Accept  json
// @Produce  json
// @Param id path string true "Shop ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Invalid ID format"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /shops/{id} [delete]
func (h *ShopHandler) DeleteShop(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.shopService.DeleteShop(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shop deleted successfully"})
}

// AddWarehouse godoc
// @Summary Add warehouse to shop
// @Description Associate a warehouse with a shop
// @Tags shops
// @Accept  json
// @Produce  json
// @Param id path string true "Shop ID"
// @Param warehouseId path string true "Warehouse ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /shops/{id}/warehouses/{warehouseId} [post]
func (h *ShopHandler) AddWarehouse(c *gin.Context) {
	shopID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID format"})
		return
	}

	warehouseID, err := primitive.ObjectIDFromHex(c.Param("warehouseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID format"})
		return
	}

	if err := h.shopService.AddWarehouseToShop(shopID, warehouseID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Warehouse added to shop successfully"})
}

// RemoveWarehouse godoc
// @Summary Remove warehouse from shop
// @Description Remove warehouse association from a shop
// @Tags shops
// @Accept  json
// @Produce  json
// @Param id path string true "Shop ID"
// @Param warehouseId path string true "Warehouse ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /shops/{id}/warehouses/{warehouseId} [delete]
func (h *ShopHandler) RemoveWarehouse(c *gin.Context) {
	shopID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID format"})
		return
	}

	warehouseID, err := primitive.ObjectIDFromHex(c.Param("warehouseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID format"})
		return
	}

	if err := h.shopService.RemoveWarehouseFromShop(shopID, warehouseID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Warehouse removed from shop successfully"})
}

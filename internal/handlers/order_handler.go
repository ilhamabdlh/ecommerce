package handlers

import (
	"net/http"

	"ecommerce/internal/models"
	"ecommerce/internal/pkg/logger"
	"ecommerce/internal/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type OrderHandler struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
	logger      *zap.Logger
}

func NewOrderHandler(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderHandler {
	return &OrderHandler{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		logger:      logger.GetLogger(),
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		h.logger.Error("Failed to bind order JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		h.logger.Error("Invalid user ID format", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	order.UserID = objID

	// Validate and reserve stock for each item
	for _, item := range order.Items {
		if err := h.productRepo.UpdateStock(item.ProductID.Hex(), item.Quantity); err != nil {
			h.logger.Error("Failed to update stock",
				zap.String("product_id", item.ProductID.Hex()),
				zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}
	}

	if err := h.orderRepo.Create(&order); err != nil {
		h.logger.Error("Failed to create order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Order created successfully", zap.String("order_id", order.ID.Hex()))
	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) ReleaseStock(c *gin.Context) {
	orderID := c.Param("id")

	if err := h.orderRepo.UpdateStatus(orderID, models.OrderStatusCanceled); err != nil {
		h.logger.Error("Failed to update order status",
			zap.String("order_id", orderID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Stock released successfully", zap.String("order_id", orderID))
	c.JSON(http.StatusOK, gin.H{"message": "Stock released successfully"})
}

// ListOrders godoc
// @Summary List all orders for authenticated user
// @Description Get all orders for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Order
// @Failure 401 {object} map[string]string
// @Router /orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	orders, err := h.orderRepo.FindByUserID(userID.(string))
	if err != nil {
		h.logger.Error("Failed to fetch orders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder godoc
// @Summary Get order by ID
// @Description Get order details by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Security Bearer
// @Success 200 {object} models.Order
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	order, err := h.orderRepo.FindByID(orderID)
	if err != nil {
		h.logger.Error("Failed to fetch order", zap.String("order_id", orderID), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

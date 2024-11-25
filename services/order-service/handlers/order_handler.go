package handlers

import (
	"net/http"

	"ecommerce/order-service/metrics"
	"ecommerce/order-service/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHandler struct {
	orderService models.OrderService
}

func NewOrderHandler(orderService models.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with the input payload
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body models.Order true "Create order"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from auth middleware and convert to ObjectID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Convert string to ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	order.UserID = userID

	if err := h.orderService.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record metrics
	metrics.OrderOperations.WithLabelValues("create", string(order.Status)).Inc()
	metrics.OrderTotalAmount.WithLabelValues(string(order.Status)).Observe(order.TotalAmount)

	c.JSON(http.StatusCreated, order)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get details of a specific order by its ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]interface{} "Invalid ID format"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Security BearerAuth
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	order, err := h.orderService.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetUserOrders godoc
// @Summary Get user's orders
// @Description Get all orders for the authenticated user
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Order
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /orders/user [get]
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Convert string to ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	orders, err := h.orderService.GetUserOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetShopOrders godoc
// @Summary Get shop's orders
// @Description Get all orders for a specific shop
// @Tags orders
// @Accept  json
// @Produce  json
// @Param shopId path string true "Shop ID"
// @Success 200 {array} models.Order
// @Failure 400 {object} map[string]interface{} "Invalid ID format"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /orders/shop/{shopId} [get]
func (h *OrderHandler) GetShopOrders(c *gin.Context) {
	shopID, err := primitive.ObjectIDFromHex(c.Param("shopId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID format"})
		return
	}

	orders, err := h.orderService.GetShopOrders(shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// ProcessOrder godoc
// @Summary Process an order
// @Description Change order status to processing
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Invalid ID format"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /orders/{id}/process [post]
func (h *OrderHandler) ProcessOrder(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.orderService.ProcessOrder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.OrderStatusTransitions.WithLabelValues(string(models.OrderStatusPending), string(models.OrderStatusProcessing)).Inc()

	c.JSON(http.StatusOK, gin.H{"message": "Order processed successfully"})
}

// CompleteOrder godoc
// @Summary Complete an order
// @Description Change order status to completed
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Invalid ID format"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /orders/{id}/complete [post]
func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.orderService.CompleteOrder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.OrderStatusTransitions.WithLabelValues(string(models.OrderStatusProcessing), string(models.OrderStatusCompleted)).Inc()

	c.JSON(http.StatusOK, gin.H{"message": "Order completed successfully"})
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Change order status to cancelled
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Invalid ID format"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.orderService.CancelOrder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.OrderStatusTransitions.WithLabelValues(string(models.OrderStatusPending), string(models.OrderStatusCancelled)).Inc()

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

package handlers

import (
	"context"
	"net/http"

	pb "github.com/ilhamabdlh/ecommerce/proto"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderClient pb.OrderServiceClient
}

func NewOrderHandler(orderClient pb.OrderServiceClient) *OrderHandler {
	return &OrderHandler{orderClient: orderClient}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req pb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	req.UserId = userID.(string)

	order, err := h.orderClient.CreateOrder(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")

	resp, err := h.orderClient.ListOrders(context.Background(), &pb.ListOrdersRequest{
		UserId: userID.(string),
		Page:   1,
		Limit:  10,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

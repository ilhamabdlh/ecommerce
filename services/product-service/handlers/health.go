package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"service": "product-service",
	})
}

func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// Add any additional readiness checks here (e.g., database connection)
	c.JSON(http.StatusOK, gin.H{
		"status":  "READY",
		"service": "product-service",
	})
}

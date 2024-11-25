package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary Check service health
// @Description Get service health status
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Service is healthy"
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"service": "order-service",
	})
}

// ReadinessCheck godoc
// @Summary Check service readiness
// @Description Get service readiness status
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Service is ready"
// @Router /ready [get]
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "READY",
		"service": "order-service",
	})
}

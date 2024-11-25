package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthHandler struct {
	db *mongo.Database
}

func NewHealthHandler(db *mongo.Database) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// @Summary Health Check
// @Description Get service health status
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Check MongoDB connection
	dbStatus := "up"
	if err := h.db.Client().Ping(ctx, nil); err != nil {
		dbStatus = "down"
	}

	response := HealthResponse{
		Status:    "UP",
		Timestamp: time.Now(),
		Services: map[string]string{
			"mongodb": dbStatus,
		},
	}

	if dbStatus == "down" {
		response.Status = "DEGRADED"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

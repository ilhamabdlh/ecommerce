package health

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthChecker struct {
	mongodb *mongo.Database
	redis   *redis.Client
	cache   map[string]time.Time
	mu      sync.RWMutex
}

func NewHealthChecker(mongodb *mongo.Database, redis *redis.Client) *HealthChecker {
	return &HealthChecker{
		mongodb: mongodb,
		redis:   redis,
		cache:   make(map[string]time.Time),
	}
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Services  map[string]string `json:"services"`
	Timestamp time.Time         `json:"timestamp"`
}

func (h *HealthChecker) CheckHealth(ctx context.Context) *HealthStatus {
	status := &HealthStatus{
		Services:  make(map[string]string),
		Timestamp: time.Now(),
	}

	// Check MongoDB
	if err := h.mongodb.Client().Ping(ctx, nil); err != nil {
		status.Services["mongodb"] = "down"
	} else {
		status.Services["mongodb"] = "up"
	}

	// Check Redis
	if err := h.redis.Ping(ctx).Err(); err != nil {
		status.Services["redis"] = "down"
	} else {
		status.Services["redis"] = "up"
	}

	// Determine overall status
	allUp := true
	for _, s := range status.Services {
		if s == "down" {
			allUp = false
			break
		}
	}

	if allUp {
		status.Status = "healthy"
	} else {
		status.Status = "unhealthy"
	}

	return status
}

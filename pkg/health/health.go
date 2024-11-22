package health

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthChecker interface {
	Check(ctx context.Context) *HealthStatus
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Services  map[string]string `json:"services"`
	Timestamp time.Time         `json:"timestamp"`
}

type HTTPHealthChecker struct {
	mongodb *mongo.Database
	redis   *redis.Client
	cache   map[string]time.Time
	mu      sync.RWMutex
}

type GRPCHealthChecker struct {
	grpc_health_v1.UnimplementedHealthServer
}

func NewHTTPHealthChecker(mongodb *mongo.Database, redis *redis.Client) *HTTPHealthChecker {
	return &HTTPHealthChecker{
		mongodb: mongodb,
		redis:   redis,
		cache:   make(map[string]time.Time),
	}
}

func NewGRPCHealthChecker() *GRPCHealthChecker {
	return &GRPCHealthChecker{}
}

func (h *HTTPHealthChecker) Check(ctx context.Context) *HealthStatus {
	status := &HealthStatus{
		Services:  make(map[string]string),
		Timestamp: time.Now(),
	}

	if err := h.mongodb.Client().Ping(ctx, nil); err != nil {
		status.Services["mongodb"] = "down"
	} else {
		status.Services["mongodb"] = "up"
	}

	if err := h.redis.Ping(ctx).Err(); err != nil {
		status.Services["redis"] = "down"
	} else {
		status.Services["redis"] = "up"
	}

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

package health

import (
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthChecker untuk HTTP health checks
type HealthChecker struct {
	mongodb *mongo.Database
	redis   *redis.Client
	cache   map[string]time.Time
	mu      sync.RWMutex
}

// GRPCHealthServer untuk gRPC health checks
type GRPCHealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
}

func NewHealthChecker(mongodb *mongo.Database, redis *redis.Client) *HealthChecker {
	return &HealthChecker{
		mongodb: mongodb,
		redis:   redis,
		cache:   make(map[string]time.Time),
	}
}

func NewGRPCHealthServer() *GRPCHealthServer {
	return &GRPCHealthServer{}
}

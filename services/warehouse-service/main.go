package main

import (
	"context"
	"ecommerce/warehouse-service/config"
	"ecommerce/warehouse-service/discovery"
	"ecommerce/warehouse-service/handlers"
	"ecommerce/warehouse-service/jobs"
	"ecommerce/warehouse-service/middleware"
	"ecommerce/warehouse-service/repository"
	"ecommerce/warehouse-service/utils"
	"net/http"
	"strconv"
	"time"

	_ "ecommerce/warehouse-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title Warehouse Service API
// @version 1.0
// @description This is a warehouse service server.
// @host localhost:8084
// @BasePath /api/v1
func setupHealthCheck(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})
}

func main() {
	// Initialize logger
	utils.InitLogger()
	defer utils.Logger.Sync()

	// Initialize circuit breaker
	utils.InitCircuitBreaker()

	// Load config
	cfg := config.LoadConfig()

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		utils.Logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(cfg.DatabaseName)

	// Initialize repository
	warehouseRepo := repository.NewMongoWarehouseRepository(db)

	// Initialize handler
	warehouseHandler := handlers.NewWarehouseHandler(warehouseRepo)
	authHandler := handlers.NewAuthHandler()

	// Setup Gin
	r := gin.Default()

	// Add middlewares
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RateLimiterMiddleware())
	r.Use(middleware.MetricsMiddleware())

	// Metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/token", authHandler.GenerateToken)
		}

		// Warehouse routes
		warehouses := api.Group("/warehouses")
		{
			warehouses.POST("/", middleware.AuthMiddleware(), middleware.RequireRoles("admin"), warehouseHandler.CreateWarehouse)
			warehouses.GET("/:id", middleware.AuthMiddleware(), warehouseHandler.GetWarehouse)
			warehouses.PUT("/:id/stock", middleware.AuthMiddleware(), middleware.RequireRoles("admin", "inventory"), warehouseHandler.UpdateStock)
			warehouses.POST("/transfer", middleware.AuthMiddleware(), middleware.RequireRoles("admin", "inventory"), warehouseHandler.TransferStock)
			warehouses.PUT("/:id/:status", middleware.AuthMiddleware(), middleware.RequireRoles("admin"), warehouseHandler.UpdateWarehouseStatus)
			warehouses.GET("/", middleware.AuthMiddleware(), warehouseHandler.GetAllWarehouses)
		}
	}

	// Setup health check
	setupHealthCheck(r)

	// Try to initialize service discovery
	if cfg.EnableServiceDiscovery {
		serviceDiscovery, err := discovery.NewServiceDiscovery(cfg.ConsulAddr)
		if err != nil {
			utils.Logger.Warnf("Failed to initialize service discovery: %v", err)
		} else {
			// Convert port string to int
			port, err := strconv.Atoi(cfg.ServerPort)
			if err != nil {
				utils.Logger.Warnf("Invalid server port: %v", err)
			} else {
				// Register service
				if err := serviceDiscovery.Register("warehouse-service", "warehouse-1", port); err != nil {
					utils.Logger.Warnf("Failed to register service: %v", err)
				} else {
					defer serviceDiscovery.Deregister("warehouse-1")
				}
			}
		}
	}

	// Initialize and start background jobs
	stockReleaseJob := jobs.NewStockReleaseJob(warehouseRepo, 5*time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stockReleaseJob.Start(ctx)

	utils.Logger.Info("Starting server on port ", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		utils.Logger.Fatalf("Failed to start server: %v", err)
	}
}

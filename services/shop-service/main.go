package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecommerce/shop-service/config"
	"ecommerce/shop-service/handlers"
	"ecommerce/shop-service/middleware"
	"ecommerce/shop-service/repository"
	"ecommerce/shop-service/services"
	"ecommerce/shop-service/utils"

	_ "ecommerce/shop-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// @title Shop Service API
// @version 1.0
// @description This is a shop service API documentation
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8083
// @BasePath /api/v1
// @schemes http
func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger, err := utils.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDB.Timeout)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoClient.Disconnect(ctx)

	db := mongoClient.Database(cfg.MongoDB.Database)
	shopCollection := db.Collection("shops")

	// Initialize repositories
	shopRepo := repository.NewMongoShopRepository(shopCollection)

	// Initialize services
	shopService := services.NewShopService(shopRepo)

	// Initialize handlers
	shopHandler := handlers.NewShopHandler(shopService)
	healthHandler := handlers.NewHealthHandler()

	// Initialize Gin router
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.ErrorHandler(logger))

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check endpoints
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/ready", healthHandler.ReadinessCheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	{
		shops := api.Group("/shops")
		{
			shops.POST("/", middleware.AuthMiddleware(cfg.JWT.Secret), shopHandler.CreateShop)
			shops.GET("/", shopHandler.ListShops)
			shops.GET("/:id", shopHandler.GetShop)
			shops.PUT("/:id", middleware.AuthMiddleware(cfg.JWT.Secret), shopHandler.UpdateShop)
			shops.DELETE("/:id", middleware.AuthMiddleware(cfg.JWT.Secret), shopHandler.DeleteShop)
			shops.POST("/:id/warehouses/:warehouseId", middleware.AuthMiddleware(cfg.JWT.Secret), shopHandler.AddWarehouse)
			shops.DELETE("/:id/warehouses/:warehouseId", middleware.AuthMiddleware(cfg.JWT.Secret), shopHandler.RemoveWarehouse)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Shutdown with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}

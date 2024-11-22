package main

import (
	"ecommerce/docs" // swagger docs
	"ecommerce/internal/config"
	"ecommerce/internal/handlers"
	"ecommerce/internal/middleware"
	"ecommerce/internal/pkg/logger"
	"ecommerce/internal/repository"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title E-Commerce API
// @version 1.0
// @description This is a sample e-commerce server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func main() {
	// Initialize logger
	logger.Init()
	log := logger.GetLogger()
	defer log.Sync()

	cfg := config.NewConfig()

	// Initialize repositories
	userRepo := repository.NewUserRepository(cfg.MongoDB)
	productRepo := repository.NewProductRepository(cfg.MongoDB)
	orderRepo := repository.NewOrderRepository(cfg.MongoDB)
	warehouseRepo := repository.NewWarehouseRepository(cfg.MongoDB)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo)
	productHandler := handlers.NewProductHandler(productRepo)
	orderHandler := handlers.NewOrderHandler(orderRepo, productRepo)
	warehouseHandler := handlers.NewWarehouseHandler(warehouseRepo, productRepo)
	healthHandler := handlers.NewHealthHandler(cfg.MongoDB)

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	r := gin.Default()

	// Swagger documentation
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Apply rate limiter to all routes
	r.Use(rateLimiter.Middleware())

	// API versioning
	v1 := r.Group("/api/v1")

	// Health check endpoint
	v1.GET("/health", healthHandler.Check)

	// Public routes
	v1.POST("/register", userHandler.Register)
	v1.POST("/login", userHandler.Login)
	v1.GET("/products", productHandler.ListProducts)

	// Protected routes
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Order endpoints
		protected.POST("/orders", orderHandler.CreateOrder)
		protected.PUT("/orders/:id/release", orderHandler.ReleaseStock)
		protected.GET("/orders", orderHandler.ListOrders)
		protected.GET("/orders/:id", orderHandler.GetOrder)

		// Product endpoints
		protected.POST("/products", productHandler.CreateProduct)
		protected.PUT("/products/:id", productHandler.UpdateProduct)
		protected.DELETE("/products/:id", productHandler.DeleteProduct)

		// Warehouse endpoints
		protected.POST("/warehouse/transfer", warehouseHandler.TransferStock)
		protected.PUT("/warehouse/:id/status", warehouseHandler.UpdateStatus)
		protected.GET("/warehouse", warehouseHandler.ListWarehouses)
		protected.GET("/warehouse/:warehouseId/products/:productId", warehouseHandler.GetProductStock)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info("Starting server", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		log.Error("Failed to start server", zap.Error(err))
		os.Exit(1)
	}
}

package main

import (
	"context"
	"log"
	"time"

	"ecommerce/user-service/config"
	_ "ecommerce/user-service/docs"
	"ecommerce/user-service/handlers"
	"ecommerce/user-service/middleware"
	"ecommerce/user-service/repository"
	"ecommerce/user-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// @title           User Service API
// @version         1.0
// @description     This is a user service API in Go using Gin framework.
// @host            localhost:8081
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Initialize logger
	utils.InitLogger()
	defer utils.Logger.Sync()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		utils.LogError("Cannot load config", zap.Error(err))
		log.Fatal("Cannot load config:", err)
	}

	// Setup rate limiter (100 requests per minute)
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	// Setup MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatal("Cannot connect to MongoDB:", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Fatal("Failed to disconnect from MongoDB:", err)
		}
	}()

	// Ping MongoDB
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatal("Cannot ping MongoDB:", err)
	}
	log.Println("Connected to MongoDB!")

	db := mongoClient.Database(cfg.MongoDB.Database)

	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	// Setup Gin
	router := gin.Default()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(rateLimiter.RateLimit())
	router.Use(middleware.MetricsMiddleware())

	// Routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// Protected routes
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/me", userHandler.GetProfile)
			users.PUT("/me", userHandler.UpdateProfile)
		}
	}

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	log.Printf("Server starting on port %s\n", cfg.Server.Port)
	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
	// test
}

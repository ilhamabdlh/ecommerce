package main

import (
	"log"
	"os"

	"github.com/ilhamabdlh/ecommerce/api-gateway/handlers"
	"github.com/ilhamabdlh/ecommerce/api-gateway/middleware"
	pbOrder "github.com/ilhamabdlh/ecommerce/proto/order"
	pbProduct "github.com/ilhamabdlh/ecommerce/proto/product"
	pbUser "github.com/ilhamabdlh/ecommerce/proto/user"
	pbWarehouse "github.com/ilhamabdlh/ecommerce/proto/warehouse"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type ServiceClients struct {
	ProductClient   pbProduct.ProductServiceClient
	OrderClient     pbOrder.OrderServiceClient
	UserClient      pbUser.UserServiceClient
	WarehouseClient pbWarehouse.WarehouseServiceClient
}

func main() {
	// Initialize gRPC connections
	clients, err := initGRPCClients()
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}

	// Initialize handlers
	productHandler := handlers.NewProductHandler(clients.ProductClient)
	orderHandler := handlers.NewOrderHandler(clients.OrderClient)
	userHandler := handlers.NewUserHandler(clients.UserClient)
	warehouseHandler := handlers.NewWarehouseHandler(clients.WarehouseClient)

	r := gin.Default()

	// Public routes
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.GET("/products", productHandler.ListProducts)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(clients.UserClient))
	{
		// Product routes
		protected.POST("/products", productHandler.CreateProduct)
		protected.PUT("/products/:id", productHandler.UpdateProduct)

		// Order routes
		protected.POST("/orders", orderHandler.CreateOrder)
		protected.GET("/orders", orderHandler.ListOrders)

		// Warehouse routes
		protected.POST("/warehouse/transfer", warehouseHandler.TransferStock)
		protected.GET("/warehouse", warehouseHandler.ListWarehouses)
	}

	r.Run(":8080")
}

func initGRPCClients() (*ServiceClients, error) {
	productConn, err := grpc.Dial(os.Getenv("PRODUCT_SERVICE_ADDR"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	orderConn, err := grpc.Dial(os.Getenv("ORDER_SERVICE_ADDR"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	userConn, err := grpc.Dial(os.Getenv("USER_SERVICE_ADDR"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	warehouseConn, err := grpc.Dial(os.Getenv("WAREHOUSE_SERVICE_ADDR"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &ServiceClients{
		ProductClient:   pbProduct.NewProductServiceClient(productConn),
		OrderClient:     pbOrder.NewOrderServiceClient(orderConn),
		UserClient:      pbUser.NewUserServiceClient(userConn),
		WarehouseClient: pbWarehouse.NewWarehouseServiceClient(warehouseConn),
	}, nil
}

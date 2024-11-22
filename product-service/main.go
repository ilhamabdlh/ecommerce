package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/ilhamabdlh/ecommerce/internal/repository"
	"github.com/ilhamabdlh/ecommerce/product-service/internal/service"
	pb "github.com/ilhamabdlh/ecommerce/proto/product"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	db := mongoClient.Database("ecommerce")
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)

	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, svc)

	log.Printf("Product service listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

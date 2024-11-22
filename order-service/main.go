package main

import (
	"log"
	"net"
	"os"

	"github.com/ilhamabdlh/ecommerce/order-service/internal/service"
	pb "github.com/ilhamabdlh/ecommerce/proto/order"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	productServiceAddr := os.Getenv("PRODUCT_SERVICE_ADDR")

	orderService := service.NewOrderService(mongoURI, productServiceAddr)
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, orderService)

	log.Printf("Order service listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

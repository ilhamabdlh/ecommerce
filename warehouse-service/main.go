package main

import (
	"log"
	"net"
	"os"

	pb "github.com/ilhamabdlh/ecommerce/proto"
	"github.com/ilhamabdlh/ecommerce/warehouse-service/internal/service"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	productServiceAddr := os.Getenv("PRODUCT_SERVICE_ADDR")

	warehouseService := service.NewWarehouseService(mongoURI, productServiceAddr)
	s := grpc.NewServer()
	pb.RegisterWarehouseServiceServer(s, warehouseService)

	log.Printf("Warehouse service listening on :50054")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

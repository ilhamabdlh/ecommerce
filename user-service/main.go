package main

import (
	"log"
	"net"
	"os"

	pb "github.com/ilhamabdlh/ecommerce/proto"
	"github.com/ilhamabdlh/ecommerce/user-service/internal/service"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	jwtSecret := os.Getenv("JWT_SECRET")

	userService := service.NewUserService(mongoURI, jwtSecret)
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userService)

	log.Printf("User service listening on :50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

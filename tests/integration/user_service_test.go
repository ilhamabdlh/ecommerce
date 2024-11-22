package integration

import (
	"context"
	"testing"
	"time"

	pb "ecommerce/proto/user"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestUserService_Register(t *testing.T) {
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	resp, err := client.Register(ctx, req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Id)
}

func TestUserService_Login(t *testing.T) {
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := client.Login(ctx, req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}

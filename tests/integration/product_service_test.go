package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/ilhamabdlh/ecommerce/proto"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestProductService_CreateProduct(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	product := &pb.Product{
		Name:  "Test Product",
		Price: 100,
		Stock: 10,
	}

	resp, err := client.CreateProduct(ctx, &pb.CreateProductRequest{Product: product})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Id)
	assert.Equal(t, product.Name, resp.Name)
}

func TestProductService_ListProducts(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.ListProducts(ctx, &pb.ListProductsRequest{
		Page:  1,
		Limit: 10,
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp.Products)
}

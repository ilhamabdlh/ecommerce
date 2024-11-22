package service

import (
	"context"
	"database/mongo"

	pb "github.com/ilhamabdlh/ecommerce/proto"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	pb.UnimplementedProductServiceServer
	db *mongo.Database
}

func NewProductService(db *mongo.Database) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	// Implementation
	return nil, status.Errorf(codes.Unimplemented, "method ListProducts not implemented")
}

func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	// Implementation
	return nil, status.Errorf(codes.Unimplemented, "method GetProduct not implemented")
}

func (s *ProductService) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
	// Implementation
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStock not implemented")
}

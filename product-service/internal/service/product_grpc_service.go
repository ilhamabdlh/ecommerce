package service

import (
	"context"

	"github.com/ilhamabdlh/ecommerce/internal/pkg/logger"
	"github.com/ilhamabdlh/ecommerce/internal/repository"
	pb "github.com/ilhamabdlh/ecommerce/proto"

	"go.uber.org/zap"
)

type ProductGRPCService struct {
	pb.UnimplementedProductServiceServer
	productRepo *repository.ProductRepository
	logger      *zap.Logger
}

func NewProductGRPCService(repo *repository.ProductRepository) *ProductGRPCService {
	return &ProductGRPCService{
		productRepo: repo,
		logger:      logger.GetLogger(),
	}
}

func (s *ProductGRPCService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := s.productRepo.List()
	if err != nil {
		s.logger.Error("Failed to list products", zap.Error(err))
		return nil, err
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:    p.ID.Hex(),
			Name:  p.Name,
			Stock: int32(p.Stock),
			Price: p.Price,
		})
	}

	return &pb.ListProductsResponse{
		Products: pbProducts,
		Total:    int32(len(products)),
	}, nil
}

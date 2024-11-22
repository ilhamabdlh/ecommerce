package service

import (
	"context"

	"github.com/ilhamabdlh/ecommerce/internal/models"
	"github.com/ilhamabdlh/ecommerce/internal/pkg/logger"
	"github.com/ilhamabdlh/ecommerce/internal/repository"
	pb "github.com/ilhamabdlh/ecommerce/proto"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type WarehouseGRPCService struct {
	pb.UnimplementedWarehouseServiceServer
	warehouseRepo *repository.WarehouseRepository
	logger        *zap.Logger
}

func NewWarehouseGRPCService(repo *repository.WarehouseRepository) *WarehouseGRPCService {
	return &WarehouseGRPCService{
		warehouseRepo: repo,
		logger:        logger.GetLogger(),
	}
}

func (s *WarehouseGRPCService) TransferStock(ctx context.Context, req *pb.TransferStockRequest) (*pb.TransferStockResponse, error) {
	transfer := &models.StockTransfer{
		ProductID:     primitive.ObjectIDFromHex(req.ProductId),
		FromWarehouse: primitive.ObjectIDFromHex(req.FromWarehouseId),
		ToWarehouse:   primitive.ObjectIDFromHex(req.ToWarehouseId),
		Quantity:      int(req.Quantity),
	}

	objID, err := primitive.ObjectIDFromHex(req.ProductId)
	if err != nil {
		return nil, err
	}
	transfer.ProductID = objID

	fromID, err := primitive.ObjectIDFromHex(req.FromWarehouseId)
	if err != nil {
		return nil, err
	}
	transfer.FromWarehouse = fromID

	toID, err := primitive.ObjectIDFromHex(req.ToWarehouseId)
	if err != nil {
		return nil, err
	}
	transfer.ToWarehouse = toID

	err := s.warehouseRepo.TransferStock(transfer)
	if err != nil {
		s.logger.Error("Failed to transfer stock", zap.Error(err))
		return &pb.TransferStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.TransferStockResponse{
		Success: true,
		Message: "Stock transferred successfully",
	}, nil
}

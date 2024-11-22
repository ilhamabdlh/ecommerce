package service

import (
	"context"
	pb "ecommerce/proto/warehouse"

	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type WarehouseService struct {
	pb.UnimplementedWarehouseServiceServer
	db            *mongo.Database
	productClient pb.ProductServiceClient
}

func NewWarehouseService(db *mongo.Database, productServiceAddr string) *WarehouseService {
	conn, err := grpc.Dial(productServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}

	return &WarehouseService{
		db:            db,
		productClient: pb.NewProductServiceClient(conn),
	}
}

func (s *WarehouseService) TransferStock(ctx context.Context, req *pb.TransferStockRequest) (*pb.TransferStockResponse, error) {
	session, err := s.db.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		fromID, err := primitive.ObjectIDFromHex(req.FromWarehouseId)
		if err != nil {
			return nil, err
		}

		toID, err := primitive.ObjectIDFromHex(req.ToWarehouseId)
		if err != nil {
			return nil, err
		}

		// Decrease stock from source warehouse
		result := s.db.Collection("products").FindOneAndUpdate(
			sessCtx,
			bson.M{
				"_id":          req.ProductId,
				"warehouse_id": fromID,
				"stock":        bson.M{"$gte": req.Quantity},
			},
			bson.M{"$inc": bson.M{"stock": -req.Quantity}},
		)
		if result.Err() != nil {
			return nil, result.Err()
		}

		// Increase stock in destination warehouse
		_, err = s.db.Collection("products").UpdateOne(
			sessCtx,
			bson.M{
				"_id":          req.ProductId,
				"warehouse_id": toID,
			},
			bson.M{"$inc": bson.M{"stock": req.Quantity}},
		)
		return nil, err
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
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

func (s *WarehouseService) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.UpdateStatusResponse, error) {
	warehouseID, err := primitive.ObjectIDFromHex(req.WarehouseId)
	if err != nil {
		return nil, err
	}

	result, err := s.db.Collection("warehouses").UpdateOne(
		ctx,
		bson.M{"_id": warehouseID},
		bson.M{"$set": bson.M{"status": req.Status}},
	)
	if err != nil {
		return &pb.UpdateStatusResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	if result.MatchedCount == 0 {
		return &pb.UpdateStatusResponse{
			Success: false,
			Message: "Warehouse not found",
		}, nil
	}

	return &pb.UpdateStatusResponse{
		Success: true,
		Message: "Warehouse status updated successfully",
	}, nil
}

func (s *WarehouseService) ListWarehouses(ctx context.Context, req *pb.ListWarehousesRequest) (*pb.ListWarehousesResponse, error) {
	cursor, err := s.db.Collection("warehouses").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var warehouses []*pb.Warehouse
	for cursor.Next(ctx) {
		var warehouse pb.Warehouse
		if err := cursor.Decode(&warehouse); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, &warehouse)
	}

	return &pb.ListWarehousesResponse{
		Warehouses: warehouses,
	}, nil
}

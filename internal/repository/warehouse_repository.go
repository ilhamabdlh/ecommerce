package repository

import (
	"context"
	"errors"
	"time"

	"ecommerce/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WarehouseRepository struct {
	db *mongo.Database
}

func NewWarehouseRepository(db *mongo.Database) *WarehouseRepository {
	return &WarehouseRepository{
		db: db,
	}
}

func (r *WarehouseRepository) Create(warehouse *models.Warehouse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	warehouse.Status = models.WarehouseStatusActive
	_, err := r.db.Collection("warehouses").InsertOne(ctx, warehouse)
	return err
}

func (r *WarehouseRepository) UpdateStatus(id string, status models.WarehouseStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	result, err := r.db.Collection("warehouses").UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("warehouse not found")
	}

	return nil
}

func (r *WarehouseRepository) TransferStock(transfer *models.StockTransfer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := r.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Kurangi stock dari warehouse asal
		_, err := r.db.Collection("products").UpdateOne(
			sessCtx,
			bson.M{
				"_id":          transfer.ProductID,
				"warehouse_id": transfer.FromWarehouse,
				"stock":        bson.M{"$gte": transfer.Quantity},
			},
			bson.M{"$inc": bson.M{"stock": -transfer.Quantity}},
		)
		if err != nil {
			return nil, err
		}

		// Tambah stock ke warehouse tujuan
		_, err = r.db.Collection("products").UpdateOne(
			sessCtx,
			bson.M{
				"_id":          transfer.ProductID,
				"warehouse_id": transfer.ToWarehouse,
			},
			bson.M{"$inc": bson.M{"stock": transfer.Quantity}},
		)
		return nil, err
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

func (r *WarehouseRepository) List() ([]models.Warehouse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.db.Collection("warehouses").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var warehouses []models.Warehouse
	if err = cursor.All(ctx, &warehouses); err != nil {
		return nil, err
	}

	return warehouses, nil
}

package mongodb

import (
	"context"
	"ecommerce/warehouse-service/internal/domain"
	"ecommerce/warehouse-service/models"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoWarehouseRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewMongoWarehouseRepository(db *mongo.Database) domain.WarehouseRepository {
	return &mongoWarehouseRepository{
		db:         db,
		collection: db.Collection("warehouses"),
	}
}

func (r *mongoWarehouseRepository) Create(ctx context.Context, warehouse *models.Warehouse) error {
	warehouse.CreatedAt = time.Now()
	warehouse.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, warehouse)
	if err != nil {
		return err
	}

	warehouse.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *mongoWarehouseRepository) GetByID(ctx context.Context, id string) (*models.Warehouse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var warehouse models.Warehouse
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&warehouse)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("warehouse not found")
		}
		return nil, err
	}

	return &warehouse, nil
}

func (r *mongoWarehouseRepository) UpdateStock(ctx context.Context, warehouseID string, productID string, quantity int) error {
	objectID, err := primitive.ObjectIDFromHex(warehouseID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$inc": bson.M{
			"stock." + productID: quantity,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("warehouse not found")
	}

	return nil
}

func (r *mongoWarehouseRepository) TransferStock(ctx context.Context, transfer *models.StockTransfer) error {
	// Start a session for transaction
	session, err := r.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) error {
		// Deduct from source warehouse
		err := r.UpdateStock(sessCtx, transfer.FromWarehouse, transfer.ProductID, -transfer.Quantity)
		if err != nil {
			return err
		}

		// Add to destination warehouse
		err = r.UpdateStock(sessCtx, transfer.ToWarehouse, transfer.ProductID, transfer.Quantity)
		if err != nil {
			return err
		}

		return nil
	}

	return mongo.WithSession(ctx, session, callback)
}

func (r *mongoWarehouseRepository) UpdateStatus(ctx context.Context, warehouseID string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(warehouseID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("warehouse not found")
	}

	return nil
}

func (r *mongoWarehouseRepository) GetAll(ctx context.Context) ([]*models.Warehouse, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var warehouses []*models.Warehouse
	if err = cursor.All(ctx, &warehouses); err != nil {
		return nil, err
	}

	return warehouses, nil
}

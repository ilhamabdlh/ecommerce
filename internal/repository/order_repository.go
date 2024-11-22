package repository

import (
	"context"
	"ecommerce/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository struct {
	db *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{
		db: db.Collection("orders"),
	}
}

func (r *OrderRepository) Create(order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = models.OrderStatusPending

	_, err := r.db.InsertOne(ctx, order)
	return err
}

func (r *OrderRepository) UpdateStatus(orderID string, status models.OrderStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err = r.db.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *OrderRepository) FindByUserID(userID string) ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.db.Find(ctx, bson.M{"user_id": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) FindByID(orderID string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, err
	}

	var order models.Order
	err = r.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

package repository

import (
	"context"
	"time"

	"ecommerce/order-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoOrderRepository struct {
	db *mongo.Collection
}

func NewMongoOrderRepository(db *mongo.Collection) models.OrderRepository {
	return &mongoOrderRepository{
		db: db,
	}
}

func (r *mongoOrderRepository) Create(order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = models.OrderStatusPending

	result, err := r.db.InsertOne(ctx, order)
	if err != nil {
		return err
	}

	order.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *mongoOrderRepository) GetByID(id primitive.ObjectID) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var order models.Order
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *mongoOrderRepository) GetByUserID(userID primitive.ObjectID) ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.db.Find(ctx, bson.M{"user_id": userID})
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

func (r *mongoOrderRepository) GetByShopID(shopID primitive.ObjectID) ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.db.Find(ctx, bson.M{"shop_id": shopID})
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

func (r *mongoOrderRepository) Update(order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order.UpdatedAt = time.Now()

	_, err := r.db.UpdateOne(
		ctx,
		bson.M{"_id": order.ID},
		bson.M{"$set": order},
	)
	return err
}

func (r *mongoOrderRepository) UpdateStatus(id primitive.ObjectID, status models.OrderStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

func (r *mongoOrderRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

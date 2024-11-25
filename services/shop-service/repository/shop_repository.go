package repository

import (
	"context"
	"time"

	"ecommerce/shop-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoShopRepository struct {
	db *mongo.Collection
}

func NewMongoShopRepository(db *mongo.Collection) models.ShopRepository {
	return &mongoShopRepository{
		db: db,
	}
}

func (r *mongoShopRepository) Create(shop *models.Shop) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shop.CreatedAt = time.Now()
	shop.UpdatedAt = time.Now()
	shop.Status = "active"
	shop.Warehouses = make([]primitive.ObjectID, 0)

	result, err := r.db.InsertOne(ctx, shop)
	if err != nil {
		return err
	}

	shop.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *mongoShopRepository) GetByID(id primitive.ObjectID) (*models.Shop, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var shop models.Shop
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&shop)
	if err != nil {
		return nil, err
	}

	return &shop, nil
}

func (r *mongoShopRepository) GetAll() ([]models.Shop, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shops []models.Shop
	if err = cursor.All(ctx, &shops); err != nil {
		return nil, err
	}

	return shops, nil
}

func (r *mongoShopRepository) Update(shop *models.Shop) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shop.UpdatedAt = time.Now()

	_, err := r.db.UpdateOne(
		ctx,
		bson.M{"_id": shop.ID},
		bson.M{"$set": shop},
	)
	return err
}

func (r *mongoShopRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mongoShopRepository) AddWarehouse(shopID, warehouseID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(
		ctx,
		bson.M{"_id": shopID},
		bson.M{
			"$addToSet": bson.M{"warehouses": warehouseID},
			"$set":      bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

func (r *mongoShopRepository) RemoveWarehouse(shopID, warehouseID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(
		ctx,
		bson.M{"_id": shopID},
		bson.M{
			"$pull": bson.M{"warehouses": warehouseID},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

func (r *mongoShopRepository) UpdateStatus(id primitive.ObjectID, status string) error {
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

package repository

import (
	"context"
	"time"

	"ecommerce/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository struct {
	db *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		db: db.Collection("products"),
	}
}

func (r *ProductRepository) Create(product *models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.db.InsertOne(ctx, product)
	return err
}

func (r *ProductRepository) List() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) UpdateStock(productID string, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$inc": bson.M{
			"stock":    -quantity,
			"reserved": quantity,
		},
	}

	opts := options.Update().SetUpsert(false)
	_, err = r.db.UpdateOne(ctx, bson.M{"_id": objID, "stock": bson.M{"$gte": quantity}}, update, opts)
	return err
}

func (r *ProductRepository) Update(productID string, product *models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":  product.Name,
			"price": product.Price,
			"stock": product.Stock,
		},
	}

	_, err = r.db.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *ProductRepository) Delete(productID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}

	_, err = r.db.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *ProductRepository) GetStockInWarehouse(productID, warehouseID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return 0, err
	}

	wID, err := primitive.ObjectIDFromHex(warehouseID)
	if err != nil {
		return 0, err
	}

	var product struct {
		Stock int `bson:"stock"`
	}
	err = r.db.FindOne(ctx, bson.M{
		"_id":          pID,
		"warehouse_id": wID,
	}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}

	return product.Stock, nil
}

package migrations

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitMigration(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create products collection if not exists
	err := db.CreateCollection(ctx, "products")
	if err != nil {
		// Ignore error if collection already exists
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}

	// Create indexes
	productsCollection := db.Collection("products")

	// Index for name field
	_, err = productsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
		Options: nil,
	})
	if err != nil {
		log.Printf("Error creating name index: %v", err)
		return err
	}

	// Index for category_id field
	_, err = productsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "category_id", Value: 1},
		},
		Options: nil,
	})
	if err != nil {
		log.Printf("Error creating category_id index: %v", err)
		return err
	}

	return nil
}

package migrations

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMigration(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create orders collection if not exists
	err := db.CreateCollection(ctx, "orders")
	if err != nil {
		// Ignore error if collection already exists
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}

	// Create indexes
	ordersCollection := db.Collection("orders")

	// Index for user_id field
	_, err = ordersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetBackground(true),
	})
	if err != nil {
		log.Printf("Error creating user_id index: %v", err)
		return err
	}

	// Index for shop_id field
	_, err = ordersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "shop_id", Value: 1},
		},
		Options: options.Index().SetBackground(true),
	})
	if err != nil {
		log.Printf("Error creating shop_id index: %v", err)
		return err
	}

	// Index for status field
	_, err = ordersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "status", Value: 1},
		},
		Options: options.Index().SetBackground(true),
	})
	if err != nil {
		log.Printf("Error creating status index: %v", err)
		return err
	}

	// Compound index for user_id and status
	_, err = ordersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "status", Value: 1},
		},
		Options: options.Index().SetBackground(true),
	})
	if err != nil {
		log.Printf("Error creating user_id_status index: %v", err)
		return err
	}

	return nil
}

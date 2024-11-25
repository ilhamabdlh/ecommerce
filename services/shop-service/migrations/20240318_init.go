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

	// Create shops collection if not exists
	err := db.CreateCollection(ctx, "shops")
	if err != nil {
		// Ignore error if collection already exists
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}

	// Create indexes
	shopsCollection := db.Collection("shops")

	// Index for name field
	_, err = shopsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
		Options: nil,
	})
	if err != nil {
		log.Printf("Error creating name index: %v", err)
		return err
	}

	// Index for status field
	_, err = shopsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "status", Value: 1},
		},
		Options: nil,
	})
	if err != nil {
		log.Printf("Error creating status index: %v", err)
		return err
	}

	// Index for location field
	_, err = shopsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "location", Value: 1},
		},
		Options: nil,
	})
	if err != nil {
		log.Printf("Error creating location index: %v", err)
		return err
	}

	return nil
}

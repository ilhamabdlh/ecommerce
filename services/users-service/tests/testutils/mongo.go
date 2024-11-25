package testutils

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupTestDB() (*mongo.Database, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create unique database name for tests
	dbName := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	db := client.Database(dbName)

	// Return cleanup function
	cleanup := func() {
		if err := db.Drop(context.Background()); err != nil {
			log.Printf("Failed to drop test database: %v", err)
		}
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}

	return db, cleanup
}

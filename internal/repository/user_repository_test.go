package repository

import (
	"context"
	"testing"
	"time"

	"ecommerce/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database("test_db")

	return db, func() {
		if err := db.Drop(context.Background()); err != nil {
			t.Errorf("Failed to drop test database: %v", err)
		}
		if err := client.Disconnect(context.Background()); err != nil {
			t.Errorf("Failed to disconnect from MongoDB: %v", err)
		}
	}
}

func TestUserRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	user := &models.User{
		Email:    "test@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	err := repo.Create(user)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	var found models.User
	err = db.Collection("users").FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&found)
	if err != nil {
		t.Errorf("Failed to find created user: %v", err)
	}

	if found.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, found.Email)
	}
}

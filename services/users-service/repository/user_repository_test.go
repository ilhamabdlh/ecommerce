package repository

import (
	"context"
	"testing"
	"time"

	"ecommerce/user-service/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database("test_db")

	return db, func() {
		if err := db.Drop(ctx); err != nil {
			t.Errorf("Failed to drop test database: %v", err)
		}
		if err := client.Disconnect(ctx); err != nil {
			t.Errorf("Failed to disconnect from MongoDB: %v", err)
		}
	}
}

func TestUserRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &models.User{
		Email:    "test@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	if user.ID.IsZero() {
		t.Error("Expected user ID to be set")
	}

	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create test user
	user := &models.User{
		Email:     "test@example.com",
		Phone:     "1234567890",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test finding user by email
	found, err := repo.FindByEmail(ctx, user.Email)
	if err != nil {
		t.Errorf("Failed to find user by email: %v", err)
	}

	if found.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, found.Email)
	}
}

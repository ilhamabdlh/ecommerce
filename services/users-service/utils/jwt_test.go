package utils

import (
	"testing"

	"ecommerce/user-service/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGenerateJWTToken(t *testing.T) {
	user := &models.User{
		ID:    primitive.NewObjectID(),
		Email: "test@example.com",
		Phone: "1234567890",
	}

	token, err := GenerateJWTToken(user)
	if err != nil {
		t.Errorf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Expected token to be non-empty")
	}
}

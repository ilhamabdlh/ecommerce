package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email" binding:"required,email"`
	Phone     string             `bson:"phone" json:"phone" binding:"required"`
	Password  string             `bson:"password" json:"password" binding:"required,min=6"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // email atau phone
	Password   string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

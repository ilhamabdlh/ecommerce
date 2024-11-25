package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Warehouse struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Location  string             `bson:"location" json:"location" validate:"required"`
	Status    string             `bson:"status" json:"status"` // "active" atau "inactive"
	Stock     map[string]int     `bson:"stock" json:"stock"`   // map productID to quantity
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type StockTransfer struct {
	ProductID     string    `bson:"product_id" json:"product_id" validate:"required"`
	Quantity      int       `bson:"quantity" json:"quantity" validate:"required,gt=0"`
	FromWarehouse string    `bson:"from_warehouse" json:"from_warehouse" validate:"required"`
	ToWarehouse   string    `bson:"to_warehouse" json:"to_warehouse" validate:"required"`
	Status        string    `bson:"status" json:"status"` // "pending", "completed", "failed"
	TransferDate  time.Time `bson:"transfer_date" json:"transfer_date"`
}

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WarehouseStatus string

const (
	WarehouseStatusActive   WarehouseStatus = "active"
	WarehouseStatusInactive WarehouseStatus = "inactive"
)

type Warehouse struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Status  WarehouseStatus    `bson:"status" json:"status"`
	Address string             `bson:"address" json:"address"`
}

type StockTransfer struct {
	ProductID     primitive.ObjectID `bson:"product_id" json:"product_id"`
	FromWarehouse primitive.ObjectID `bson:"from_warehouse" json:"from_warehouse"`
	ToWarehouse   primitive.ObjectID `bson:"to_warehouse" json:"to_warehouse"`
	Quantity      int                `bson:"quantity" json:"quantity"`
}

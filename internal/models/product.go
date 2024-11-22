package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Stock    int                `bson:"stock" json:"stock"`
	Price    float64            `bson:"price" json:"price"`
	Reserved int                `bson:"reserved" json:"reserved"`
}

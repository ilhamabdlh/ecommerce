package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Stock       int                `bson:"stock" json:"stock"`
	CategoryID  primitive.ObjectID `bson:"category_id" json:"category_id"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type ProductRepository interface {
	Create(product *Product) error
	GetByID(id primitive.ObjectID) (*Product, error)
	GetAll() ([]Product, error)
	Update(product *Product) error
	Delete(id primitive.ObjectID) error
	UpdateStock(id primitive.ObjectID, quantity int) error
}

type ProductService interface {
	CreateProduct(product *Product) error
	GetProduct(id primitive.ObjectID) (*Product, error)
	ListProducts() ([]Product, error)
	UpdateProduct(product *Product) error
	DeleteProduct(id primitive.ObjectID) error
	UpdateProductStock(id primitive.ObjectID, quantity int) error
}

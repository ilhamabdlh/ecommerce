package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shop struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	Location    string               `bson:"location" json:"location"`
	Status      string               `bson:"status" json:"status"` // active, inactive
	Warehouses  []primitive.ObjectID `bson:"warehouses" json:"warehouses"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updated_at"`
}

func NewShop() *Shop {
	return &Shop{
		Warehouses: make([]primitive.ObjectID, 0),
		Status:     "active",
	}
}

type ShopRepository interface {
	Create(shop *Shop) error
	GetByID(id primitive.ObjectID) (*Shop, error)
	GetAll() ([]Shop, error)
	Update(shop *Shop) error
	Delete(id primitive.ObjectID) error
	AddWarehouse(shopID, warehouseID primitive.ObjectID) error
	RemoveWarehouse(shopID, warehouseID primitive.ObjectID) error
	UpdateStatus(id primitive.ObjectID, status string) error
}

type ShopService interface {
	CreateShop(shop *Shop) error
	GetShop(id primitive.ObjectID) (*Shop, error)
	ListShops() ([]Shop, error)
	UpdateShop(shop *Shop) error
	DeleteShop(id primitive.ObjectID) error
	AddWarehouseToShop(shopID, warehouseID primitive.ObjectID) error
	RemoveWarehouseFromShop(shopID, warehouseID primitive.ObjectID) error
	ActivateShop(id primitive.ObjectID) error
	DeactivateShop(id primitive.ObjectID) error
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type OrderItem struct {
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	Price     float64            `bson:"price" json:"price"`
}

type Order struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID  `bson:"user_id" json:"user_id"`
	ShopID      primitive.ObjectID  `bson:"shop_id" json:"shop_id"`
	Items       []OrderItem         `bson:"items" json:"items"`
	TotalAmount float64             `bson:"total_amount" json:"total_amount"`
	Status      OrderStatus         `bson:"status" json:"status"`
	PaymentID   *primitive.ObjectID `bson:"payment_id,omitempty" json:"payment_id,omitempty"`
	CreatedAt   time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at" json:"updated_at"`
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id primitive.ObjectID) (*Order, error)
	GetByUserID(userID primitive.ObjectID) ([]Order, error)
	GetByShopID(shopID primitive.ObjectID) ([]Order, error)
	Update(order *Order) error
	UpdateStatus(id primitive.ObjectID, status OrderStatus) error
	Delete(id primitive.ObjectID) error
}

type OrderService interface {
	CreateOrder(order *Order) error
	GetOrder(id primitive.ObjectID) (*Order, error)
	GetUserOrders(userID primitive.ObjectID) ([]Order, error)
	GetShopOrders(shopID primitive.ObjectID) ([]Order, error)
	UpdateOrder(order *Order) error
	ProcessOrder(id primitive.ObjectID) error
	CompleteOrder(id primitive.ObjectID) error
	CancelOrder(id primitive.ObjectID) error
}

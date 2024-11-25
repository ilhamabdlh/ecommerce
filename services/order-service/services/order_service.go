package services

import (
	"errors"

	"ecommerce/order-service/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type orderService struct {
	orderRepo models.OrderRepository
}

func NewOrderService(orderRepo models.OrderRepository) models.OrderService {
	return &orderService{
		orderRepo: orderRepo,
	}
}

func (s *orderService) CreateOrder(order *models.Order) error {
	if order.UserID.IsZero() {
		return errors.New("user ID is required")
	}
	if order.ShopID.IsZero() {
		return errors.New("shop ID is required")
	}
	if len(order.Items) == 0 {
		return errors.New("order must have at least one item")
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range order.Items {
		if item.Quantity <= 0 {
			return errors.New("item quantity must be greater than 0")
		}
		if item.Price <= 0 {
			return errors.New("item price must be greater than 0")
		}
		totalAmount += float64(item.Quantity) * item.Price
	}
	order.TotalAmount = totalAmount
	order.Status = models.OrderStatusPending

	return s.orderRepo.Create(order)
}

func (s *orderService) GetOrder(id primitive.ObjectID) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *orderService) GetUserOrders(userID primitive.ObjectID) ([]models.Order, error) {
	return s.orderRepo.GetByUserID(userID)
}

func (s *orderService) GetShopOrders(shopID primitive.ObjectID) ([]models.Order, error) {
	return s.orderRepo.GetByShopID(shopID)
}

func (s *orderService) UpdateOrder(order *models.Order) error {
	if order.ID.IsZero() {
		return errors.New("order ID is required")
	}
	if order.UserID.IsZero() {
		return errors.New("user ID is required")
	}
	if order.ShopID.IsZero() {
		return errors.New("shop ID is required")
	}
	if len(order.Items) == 0 {
		return errors.New("order must have at least one item")
	}

	// Recalculate total amount
	var totalAmount float64
	for _, item := range order.Items {
		if item.Quantity <= 0 {
			return errors.New("item quantity must be greater than 0")
		}
		if item.Price <= 0 {
			return errors.New("item price must be greater than 0")
		}
		totalAmount += float64(item.Quantity) * item.Price
	}
	order.TotalAmount = totalAmount

	return s.orderRepo.Update(order)
}

func (s *orderService) ProcessOrder(id primitive.ObjectID) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}

	if order.Status != models.OrderStatusPending {
		return errors.New("order must be in pending status to process")
	}

	return s.orderRepo.UpdateStatus(id, models.OrderStatusProcessing)
}

func (s *orderService) CompleteOrder(id primitive.ObjectID) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}

	if order.Status != models.OrderStatusProcessing {
		return errors.New("order must be in processing status to complete")
	}

	return s.orderRepo.UpdateStatus(id, models.OrderStatusCompleted)
}

func (s *orderService) CancelOrder(id primitive.ObjectID) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}

	if order.Status == models.OrderStatusCompleted {
		return errors.New("completed order cannot be cancelled")
	}

	return s.orderRepo.UpdateStatus(id, models.OrderStatusCancelled)
}

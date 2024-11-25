package tests

import (
	"testing"
	"time"

	"ecommerce/order-service/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(id primitive.ObjectID) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByUserID(userID primitive.ObjectID) ([]models.Order, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByShopID(shopID primitive.ObjectID) ([]models.Order, error) {
	args := m.Called(shopID)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateStatus(id primitive.ObjectID, status models.OrderStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	order := &models.Order{
		UserID: primitive.NewObjectID(),
		ShopID: primitive.NewObjectID(),
		Items: []models.OrderItem{
			{
				ProductID: primitive.NewObjectID(),
				Quantity:  2,
				Price:     100.0,
			},
		},
		Status:    models.OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("Create", order).Return(nil)

	assert.NoError(t, mockRepo.Create(order))
	mockRepo.AssertExpectations(t)
}

func TestGetOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	id := primitive.NewObjectID()
	expectedOrder := &models.Order{
		ID:          id,
		UserID:      primitive.NewObjectID(),
		ShopID:      primitive.NewObjectID(),
		Status:      models.OrderStatusPending,
		TotalAmount: 200.0,
	}

	mockRepo.On("GetByID", id).Return(expectedOrder, nil)

	order, err := mockRepo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
	mockRepo.AssertExpectations(t)
}

func TestGetUserOrders(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	userID := primitive.NewObjectID()
	expectedOrders := []models.Order{
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Status:      models.OrderStatusPending,
			TotalAmount: 100.0,
		},
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Status:      models.OrderStatusCompleted,
			TotalAmount: 200.0,
		},
	}

	mockRepo.On("GetByUserID", userID).Return(expectedOrders, nil)

	orders, err := mockRepo.GetByUserID(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, orders)
	mockRepo.AssertExpectations(t)
}

func TestUpdateOrderStatus(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	id := primitive.NewObjectID()
	newStatus := models.OrderStatusProcessing

	mockRepo.On("UpdateStatus", id, newStatus).Return(nil)

	err := mockRepo.UpdateStatus(id, newStatus)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

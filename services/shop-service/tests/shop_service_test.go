package tests

import (
	"testing"
	"time"

	"ecommerce/shop-service/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockShopRepository struct {
	mock.Mock
}

func (m *MockShopRepository) Create(shop *models.Shop) error {
	args := m.Called(shop)
	return args.Error(0)
}

func (m *MockShopRepository) GetByID(id primitive.ObjectID) (*models.Shop, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Shop), args.Error(1)
}

func (m *MockShopRepository) GetAll() ([]models.Shop, error) {
	args := m.Called()
	return args.Get(0).([]models.Shop), args.Error(1)
}

func (m *MockShopRepository) Update(shop *models.Shop) error {
	args := m.Called(shop)
	return args.Error(0)
}

func (m *MockShopRepository) Delete(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockShopRepository) AddWarehouse(shopID, warehouseID primitive.ObjectID) error {
	args := m.Called(shopID, warehouseID)
	return args.Error(0)
}

func (m *MockShopRepository) RemoveWarehouse(shopID, warehouseID primitive.ObjectID) error {
	args := m.Called(shopID, warehouseID)
	return args.Error(0)
}

func (m *MockShopRepository) UpdateStatus(id primitive.ObjectID, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func TestCreateShop(t *testing.T) {
	mockRepo := new(MockShopRepository)
	shop := &models.Shop{
		Name:        "Test Shop",
		Description: "Test Description",
		Location:    "Test Location",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("Create", shop).Return(nil)

	assert.NoError(t, mockRepo.Create(shop))
	mockRepo.AssertExpectations(t)
}

func TestGetShop(t *testing.T) {
	mockRepo := new(MockShopRepository)
	id := primitive.NewObjectID()
	expectedShop := &models.Shop{
		ID:          id,
		Name:        "Test Shop",
		Description: "Test Description",
		Location:    "Test Location",
		Status:      "active",
	}

	mockRepo.On("GetByID", id).Return(expectedShop, nil)

	shop, err := mockRepo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, expectedShop, shop)
	mockRepo.AssertExpectations(t)
}

func TestAddWarehouse(t *testing.T) {
	mockRepo := new(MockShopRepository)
	shopID := primitive.NewObjectID()
	warehouseID := primitive.NewObjectID()

	mockRepo.On("AddWarehouse", shopID, warehouseID).Return(nil)

	err := mockRepo.AddWarehouse(shopID, warehouseID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

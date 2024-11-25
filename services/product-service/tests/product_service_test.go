package tests

import (
	"testing"
	"time"

	"ecommerce/product-service/models"
	"ecommerce/product-service/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(id primitive.ObjectID) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetAll() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateStock(id primitive.ObjectID, quantity int) error {
	args := m.Called(id, quantity)
	return args.Error(0)
}

func TestCreateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := services.NewProductService(mockRepo)

	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		Stock:       10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("Create", product).Return(nil)

	err := productService.CreateProduct(product)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := services.NewProductService(mockRepo)

	id := primitive.NewObjectID()
	expectedProduct := &models.Product{
		ID:          id,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		Stock:       10,
	}

	mockRepo.On("GetByID", id).Return(expectedProduct, nil)

	product, err := productService.GetProduct(id)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProductStock(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productService := services.NewProductService(mockRepo)

	id := primitive.NewObjectID()
	existingProduct := &models.Product{
		ID:    id,
		Stock: 10,
	}

	mockRepo.On("GetByID", id).Return(existingProduct, nil)
	mockRepo.On("UpdateStock", id, 5).Return(nil)

	err := productService.UpdateProductStock(id, 5)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

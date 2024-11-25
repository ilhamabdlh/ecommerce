package tests

import (
	"context"
	"ecommerce/warehouse-service/models"
	"ecommerce/warehouse-service/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockWarehouseRepo struct {
	mock.Mock
	repository.WarehouseRepository
}

func (m *MockWarehouseRepo) Create(ctx context.Context, warehouse *models.Warehouse) error {
	args := m.Called(ctx, warehouse)
	return args.Error(0)
}

func (m *MockWarehouseRepo) GetByID(ctx context.Context, id string) (*models.Warehouse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Warehouse), args.Error(1)
}

func TestWarehouseOperations(t *testing.T) {
	t.Run("CreateWarehouse", func(t *testing.T) {
		mockRepo := new(MockWarehouseRepo)
		warehouse := &models.Warehouse{
			ID:       primitive.NewObjectID(),
			Name:     "Test Warehouse",
			Location: "Test Location",
			Status:   "active",
			Stock:    make(map[string]int),
		}

		mockRepo.On("Create", mock.Anything, warehouse).Return(nil)

		err := mockRepo.Create(context.Background(), warehouse)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetWarehouse", func(t *testing.T) {
		mockRepo := new(MockWarehouseRepo)
		warehouse := &models.Warehouse{
			ID:       primitive.NewObjectID(),
			Name:     "Test Warehouse",
			Location: "Test Location",
			Status:   "active",
			Stock:    make(map[string]int),
		}

		mockRepo.On("GetByID", mock.Anything, warehouse.ID.Hex()).Return(warehouse, nil)

		result, err := mockRepo.GetByID(context.Background(), warehouse.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, warehouse, result)
		mockRepo.AssertExpectations(t)
	})
}

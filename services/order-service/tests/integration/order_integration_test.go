package integration

import (
	"context"
	"testing"
	"time"

	"ecommerce/order-service/config"
	"ecommerce/order-service/models"
	"ecommerce/order-service/repository"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Collection, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := config.LoadConfig()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	assert.NoError(t, err)

	db := client.Database("test_db")
	collection := db.Collection("orders")

	return collection, func() {
		collection.Drop(ctx)
		client.Disconnect(ctx)
	}
}

func TestOrderRepositoryIntegration(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewMongoOrderRepository(collection)

	t.Run("Create and Get Order", func(t *testing.T) {
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
			TotalAmount: 200.0,
			Status:      models.OrderStatusPending,
		}

		err := repo.Create(order)
		assert.NoError(t, err)
		assert.NotEmpty(t, order.ID)

		retrieved, err := repo.GetByID(order.ID)
		assert.NoError(t, err)
		assert.Equal(t, order.UserID, retrieved.UserID)
		assert.Equal(t, order.TotalAmount, retrieved.TotalAmount)
	})

	t.Run("Get User Orders", func(t *testing.T) {
		userID := primitive.NewObjectID()
		order1 := &models.Order{
			UserID:      userID,
			ShopID:      primitive.NewObjectID(),
			TotalAmount: 100.0,
			Status:      models.OrderStatusPending,
		}
		order2 := &models.Order{
			UserID:      userID,
			ShopID:      primitive.NewObjectID(),
			TotalAmount: 200.0,
			Status:      models.OrderStatusCompleted,
		}

		err := repo.Create(order1)
		assert.NoError(t, err)
		err = repo.Create(order2)
		assert.NoError(t, err)

		orders, err := repo.GetByUserID(userID)
		assert.NoError(t, err)
		assert.Len(t, orders, 2)
	})

	t.Run("Update Order Status", func(t *testing.T) {
		order := &models.Order{
			UserID:      primitive.NewObjectID(),
			ShopID:      primitive.NewObjectID(),
			Status:      models.OrderStatusPending,
			TotalAmount: 150.0,
		}

		err := repo.Create(order)
		assert.NoError(t, err)

		err = repo.UpdateStatus(order.ID, models.OrderStatusProcessing)
		assert.NoError(t, err)

		updated, err := repo.GetByID(order.ID)
		assert.NoError(t, err)
		assert.Equal(t, models.OrderStatusProcessing, updated.Status)
	})
}

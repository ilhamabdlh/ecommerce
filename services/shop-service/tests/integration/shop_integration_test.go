package integration

import (
	"context"
	"testing"
	"time"

	"ecommerce/shop-service/config"
	"ecommerce/shop-service/models"
	"ecommerce/shop-service/repository"

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
	collection := db.Collection("shops")

	return collection, func() {
		collection.Drop(ctx)
		client.Disconnect(ctx)
	}
}

func TestShopRepositoryIntegration(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewMongoShopRepository(collection)

	t.Run("Create and Get Shop", func(t *testing.T) {
		shop := &models.Shop{
			Name:        "Integration Test Shop",
			Description: "Test Description",
			Location:    "Test Location",
			Status:      "active",
		}

		err := repo.Create(shop)
		assert.NoError(t, err)
		assert.NotEmpty(t, shop.ID)

		retrieved, err := repo.GetByID(shop.ID)
		assert.NoError(t, err)
		assert.Equal(t, shop.Name, retrieved.Name)
		assert.Equal(t, shop.Location, retrieved.Location)
	})

	t.Run("Update Shop", func(t *testing.T) {
		shop := &models.Shop{
			Name:        "Shop to Update",
			Description: "Original Description",
			Location:    "Original Location",
			Status:      "active",
		}

		err := repo.Create(shop)
		assert.NoError(t, err)

		shop.Description = "Updated Description"
		shop.Location = "Updated Location"
		err = repo.Update(shop)
		assert.NoError(t, err)

		updated, err := repo.GetByID(shop.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Description", updated.Description)
		assert.Equal(t, "Updated Location", updated.Location)
	})

	t.Run("Delete Shop", func(t *testing.T) {
		shop := &models.Shop{
			Name:     "Shop to Delete",
			Location: "Test Location",
			Status:   "active",
		}

		err := repo.Create(shop)
		assert.NoError(t, err)

		err = repo.Delete(shop.ID)
		assert.NoError(t, err)

		_, err = repo.GetByID(shop.ID)
		assert.Error(t, err)
	})

	t.Run("Add and Remove Warehouse", func(t *testing.T) {
		shop := &models.Shop{
			Name:       "Shop with Warehouse",
			Location:   "Test Location",
			Status:     "active",
			Warehouses: make([]primitive.ObjectID, 0),
		}

		err := repo.Create(shop)
		assert.NoError(t, err)

		warehouseID := primitive.NewObjectID()
		err = repo.AddWarehouse(shop.ID, warehouseID)
		assert.NoError(t, err)

		updated, err := repo.GetByID(shop.ID)
		assert.NoError(t, err)
		assert.Contains(t, updated.Warehouses, warehouseID)

		err = repo.RemoveWarehouse(shop.ID, warehouseID)
		assert.NoError(t, err)

		updated, err = repo.GetByID(shop.ID)
		assert.NoError(t, err)
		assert.NotContains(t, updated.Warehouses, warehouseID)
	})

	t.Run("Update Shop Status", func(t *testing.T) {
		shop := &models.Shop{
			Name:     "Shop Status Test",
			Location: "Test Location",
			Status:   "active",
		}

		err := repo.Create(shop)
		assert.NoError(t, err)

		err = repo.UpdateStatus(shop.ID, "inactive")
		assert.NoError(t, err)

		updated, err := repo.GetByID(shop.ID)
		assert.NoError(t, err)
		assert.Equal(t, "inactive", updated.Status)
	})
}

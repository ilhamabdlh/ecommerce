package integration

import (
	"context"
	"testing"
	"time"

	"ecommerce/product-service/config"
	"ecommerce/product-service/models"
	"ecommerce/product-service/repository"

	"github.com/stretchr/testify/assert"
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
	collection := db.Collection("products")

	return collection, func() {
		collection.Drop(ctx)
		client.Disconnect(ctx)
	}
}

func TestProductRepositoryIntegration(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewMongoProductRepository(collection)

	t.Run("Create and Get Product", func(t *testing.T) {
		product := &models.Product{
			Name:        "Integration Test Product",
			Description: "Test Description",
			Price:       100.0,
			Stock:       10,
		}

		err := repo.Create(product)
		assert.NoError(t, err)
		assert.NotEmpty(t, product.ID)

		retrieved, err := repo.GetByID(product.ID)
		assert.NoError(t, err)
		assert.Equal(t, product.Name, retrieved.Name)
	})

	t.Run("Update Product", func(t *testing.T) {
		product := &models.Product{
			Name:        "Product to Update",
			Description: "Original Description",
			Price:       100.0,
			Stock:       10,
		}

		err := repo.Create(product)
		assert.NoError(t, err)

		product.Description = "Updated Description"
		err = repo.Update(product)
		assert.NoError(t, err)

		updated, err := repo.GetByID(product.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Description", updated.Description)
	})

	t.Run("Delete Product", func(t *testing.T) {
		product := &models.Product{
			Name:  "Product to Delete",
			Price: 100.0,
			Stock: 10,
		}

		err := repo.Create(product)
		assert.NoError(t, err)

		err = repo.Delete(product.ID)
		assert.NoError(t, err)

		_, err = repo.GetByID(product.ID)
		assert.Error(t, err)
	})
}

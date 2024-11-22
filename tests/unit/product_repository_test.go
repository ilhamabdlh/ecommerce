package unit

import (
	"context"
	"testing"
	"time"

	"ecommerce/internal/models"
	"ecommerce/internal/repository"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)

	return client.Database("ecommerce_test")
}

func TestProductRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProductRepository(db)

	product := &models.Product{
		Name:  "Test Product",
		Price: 100,
		Stock: 10,
	}

	err := repo.Create(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
}

func TestProductRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProductRepository(db)

	// Create test products
	products := []models.Product{
		{Name: "Product 1", Price: 100, Stock: 10},
		{Name: "Product 2", Price: 200, Stock: 20},
	}

	for _, p := range products {
		err := repo.Create(&p)
		assert.NoError(t, err)
	}

	// Test listing
	result, err := repo.List()
	assert.NoError(t, err)
	assert.Len(t, result, len(products))
}

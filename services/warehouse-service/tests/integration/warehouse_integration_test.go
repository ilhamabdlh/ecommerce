package integration

import (
	"context"
	"ecommerce/warehouse-service/models"
	"ecommerce/warehouse-service/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)

	db := client.Database("warehouse_test_db")

	return db, func() {
		err := db.Drop(ctx)
		assert.NoError(t, err)
		err = client.Disconnect(ctx)
		assert.NoError(t, err)
	}
}

func TestWarehouseIntegration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewMongoWarehouseRepository(db)

	t.Run("Create and Get Warehouse", func(t *testing.T) {
		warehouse := &models.Warehouse{
			Name:     "Test Warehouse",
			Location: "Test Location",
			Status:   "active",
			Stock:    make(map[string]int),
		}

		// Create warehouse
		err := repo.Create(context.Background(), warehouse)
		assert.NoError(t, err)
		assert.NotEmpty(t, warehouse.ID)

		// Get warehouse
		retrieved, err := repo.GetByID(context.Background(), warehouse.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, warehouse.Name, retrieved.Name)
		assert.Equal(t, warehouse.Location, retrieved.Location)
	})

	t.Run("Stock Operations", func(t *testing.T) {
		warehouse := &models.Warehouse{
			Name:     "Stock Test Warehouse",
			Location: "Test Location",
			Status:   "active",
			Stock:    make(map[string]int),
		}

		// Create warehouse
		err := repo.Create(context.Background(), warehouse)
		assert.NoError(t, err)

		// Update stock
		productID := "test-product"
		err = repo.UpdateStock(context.Background(), warehouse.ID.Hex(), productID, 10)
		assert.NoError(t, err)

		// Verify stock update
		updated, err := repo.GetByID(context.Background(), warehouse.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, 10, updated.Stock[productID])
	})

	t.Run("Transfer Stock", func(t *testing.T) {
		// Create source warehouse
		sourceWarehouse := &models.Warehouse{
			Name:     "Source Warehouse",
			Location: "Location A",
			Status:   "active",
			Stock:    map[string]int{"product1": 100},
		}
		err := repo.Create(context.Background(), sourceWarehouse)
		assert.NoError(t, err)

		// Create destination warehouse
		destWarehouse := &models.Warehouse{
			Name:     "Destination Warehouse",
			Location: "Location B",
			Status:   "active",
			Stock:    make(map[string]int),
		}
		err = repo.Create(context.Background(), destWarehouse)
		assert.NoError(t, err)

		// Perform transfer
		transfer := &models.StockTransfer{
			ProductID:     "product1",
			Quantity:      50,
			FromWarehouse: sourceWarehouse.ID.Hex(),
			ToWarehouse:   destWarehouse.ID.Hex(),
			Status:        "pending",
			TransferDate:  time.Now(),
		}

		err = repo.TransferStock(context.Background(), transfer)
		assert.NoError(t, err)

		// Verify transfer results
		source, err := repo.GetByID(context.Background(), sourceWarehouse.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, 50, source.Stock["product1"])

		dest, err := repo.GetByID(context.Background(), destWarehouse.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, 50, dest.Stock["product1"])
	})
}

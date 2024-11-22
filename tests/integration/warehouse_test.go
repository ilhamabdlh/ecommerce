package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilhamabdlh/ecommerce/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createTestWarehouse(t *testing.T, r *gin.Engine, token string, name string) *models.Warehouse {
	warehouse := models.Warehouse{
		Name:    name,
		Status:  models.WarehouseStatusActive,
		Address: "Test Address",
	}

	jsonValue, _ := json.Marshal(warehouse)
	req, _ := http.NewRequest("POST", "/api/v1/warehouse", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Warehouse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	return &response
}

func createTestProductInWarehouse(t *testing.T, r *gin.Engine, token string, warehouseID primitive.ObjectID) *models.Product {
	product := &models.Product{
		Name:  "Test Product",
		Stock: 10,
		Price: 100,
	}

	jsonValue, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Warehouse-ID", warehouseID.Hex())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdProduct models.Product
	err := json.Unmarshal(w.Body.Bytes(), &createdProduct)
	assert.NoError(t, err)
	return &createdProduct
}

func TestWarehouseTransfer(t *testing.T) {
	r := setupTestRouter()
	token := loginTestUser(t, r)

	// Create source warehouse
	sourceWarehouse := createTestWarehouse(t, r, token, "Source Warehouse")
	// Create destination warehouse
	destWarehouse := createTestWarehouse(t, r, token, "Destination Warehouse")

	// Create product with initial stock in source warehouse
	product := createTestProductInWarehouse(t, r, token, sourceWarehouse.ID)

	// Test stock transfer
	transfer := models.StockTransfer{
		ProductID:     product.ID,
		FromWarehouse: sourceWarehouse.ID,
		ToWarehouse:   destWarehouse.ID,
		Quantity:      5,
	}

	jsonValue, _ := json.Marshal(transfer)
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/transfer", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify stock in both warehouses
	verifyWarehouseStock(t, r, token, product.ID, sourceWarehouse.ID, 5)
	verifyWarehouseStock(t, r, token, product.ID, destWarehouse.ID, 5)
}

func verifyWarehouseStock(t *testing.T, r *gin.Engine, token string, productID, warehouseID primitive.ObjectID, expectedStock int) {
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/"+warehouseID.Hex()+"/products/"+productID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Stock int `json:"stock"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedStock, response.Stock)
}

func TestWarehouseStatus(t *testing.T) {
	r := setupTestRouter()
	token := loginTestUser(t, r)

	// Create warehouse
	warehouse := createTestWarehouse(t, r, token, "Test Warehouse")

	// Update warehouse status
	updateStatus := struct {
		Status models.WarehouseStatus `json:"status"`
	}{
		Status: models.WarehouseStatusInactive,
	}

	jsonValue, _ := json.Marshal(updateStatus)
	req, _ := http.NewRequest("PUT", "/api/v1/warehouse/"+warehouse.ID.Hex()+"/status", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify warehouse status
	req, _ = http.NewRequest("GET", "/api/v1/warehouse/"+warehouse.ID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedWarehouse models.Warehouse
	err := json.Unmarshal(w.Body.Bytes(), &updatedWarehouse)
	assert.NoError(t, err)
	assert.Equal(t, models.WarehouseStatusInactive, updatedWarehouse.Status)
}

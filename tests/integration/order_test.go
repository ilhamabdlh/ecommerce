package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecommerce/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateOrder(t *testing.T) {
	r := setupTestRouter()
	token := loginTestUser(t, r)

	productID := createTestProduct(t, r, token)

	order := models.Order{
		Items: []models.OrderItem{
			{
				ProductID: productID,
				Quantity:  1,
				Price:     100,
			},
		},
	}

	jsonValue, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Order
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, models.OrderStatusPending, response.Status)
}

func createTestProduct(t *testing.T, r *gin.Engine, token string) primitive.ObjectID {
	product := models.Product{
		Name:  "Test Product",
		Stock: 10,
		Price: 100,
	}

	jsonValue, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Product
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response.ID
}

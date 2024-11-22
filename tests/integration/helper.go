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
)

func loginTestUser(t *testing.T, r *gin.Engine) string {
	user := models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Register user first
	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Login
	req, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	return response["token"].(string)
}

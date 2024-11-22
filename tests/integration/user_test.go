package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecommerce/internal/config"
	"ecommerce/internal/handlers"
	"ecommerce/internal/models"
	"ecommerce/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	cfg := config.NewConfig()
	userRepo := repository.NewUserRepository(cfg.MongoDB)
	userHandler := handlers.NewUserHandler(userRepo)

	r := gin.Default()
	r.POST("/register", userHandler.Register)
	return r
}

func TestUserRegistration(t *testing.T) {
	r := setupTestRouter()

	user := models.User{
		Email:    "test@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User created successfully", response["message"])
}

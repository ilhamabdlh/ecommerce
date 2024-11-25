package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ecommerce/user-service/handlers"
	"ecommerce/user-service/middleware"
	"ecommerce/user-service/models"
	"ecommerce/user-service/repository"
	"ecommerce/user-service/tests/testutils"
	"ecommerce/user-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestServer(t *testing.T) (*gin.Engine, *mongo.Database, func()) {
	// Setup logger
	utils.Logger = testutils.NewTestLogger(t)

	// Setup MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	dbName := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	db := client.Database(dbName)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup repositories and handlers
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	// Setup routes
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/me", userHandler.GetProfile)
			users.PUT("/me", userHandler.UpdateProfile)
		}
	}

	// Cleanup function
	cleanup := func() {
		if err := db.Drop(context.Background()); err != nil {
			t.Errorf("Failed to drop test database: %v", err)
		}
		if err := client.Disconnect(context.Background()); err != nil {
			t.Errorf("Failed to disconnect from MongoDB: %v", err)
		}
	}

	return router, db, cleanup
}

func TestUserRegistrationAndLogin(t *testing.T) {
	router, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Test Registration
	registerPayload := models.User{
		Email:    "test@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	body, _ := json.Marshal(registerPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var registerResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &registerResponse)
	assert.NoError(t, err)
	assert.Contains(t, registerResponse, "user_id")

	// Test Login
	loginPayload := models.LoginRequest{
		Identifier: registerPayload.Email,
		Password:   registerPayload.Password,
	}

	body, _ = json.Marshal(loginPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var loginResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.NoError(t, err)
	assert.Contains(t, loginResponse, "token")
}

func TestProtectedEndpoints(t *testing.T) {
	router, _, cleanup := setupTestServer(t)
	defer cleanup()

	// First register and login to get token
	user := models.User{
		Email:    "test@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	// Register
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Login to get token
	loginPayload := models.LoginRequest{
		Identifier: user.Email,
		Password:   user.Password,
	}

	body, _ = json.Marshal(loginPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	// Test GetProfile endpoint
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var profile models.User
	err := json.Unmarshal(w.Body.Bytes(), &profile)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, profile.Email)
}

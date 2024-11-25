package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecommerce/user-service/models"
	"ecommerce/user-service/repository"
	"ecommerce/user-service/tests/testutils"
	"ecommerce/user-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func setupTest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	utils.Logger = testutils.NewTestLogger(t)
}

// Mock repository
type MockUserRepository struct {
	mock.Mock
}

// Implement UserRepositoryInterface
var _ repository.UserRepositoryInterface = (*MockUserRepository)(nil)

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id string, user *models.User) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

func TestUserHandler_Register(t *testing.T) {
	setupTest(t)

	tests := []struct {
		name         string
		input        models.User
		mockBehavior func(*MockUserRepository)
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name: "Success",
			input: models.User{
				Email:    "test@example.com",
				Phone:    "1234567890",
				Password: "password123",
			},
			mockBehavior: func(repo *MockUserRepository) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "User created successfully",
			},
		},
		{
			name: "Invalid Email",
			input: models.User{
				Email:    "invalid-email",
				Phone:    "1234567890",
				Password: "password123",
			},
			mockBehavior: func(repo *MockUserRepository) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing Required Fields",
			input: models.User{
				Email: "test@example.com",
			},
			mockBehavior: func(repo *MockUserRepository) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			tt.mockBehavior(mockRepo)
			handler := NewUserHandler(mockRepo)

			// Create request
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Test
			handler.Register(c)

			// Assert status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// If expected body is set, verify response
			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				for k, v := range tt.expectedBody {
					assert.Equal(t, v, response[k])
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	setupTest(t)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &models.User{
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	tests := []struct {
		name         string
		input        models.LoginRequest
		mockBehavior func(*MockUserRepository)
		expectedCode int
	}{
		{
			name: "Success",
			input: models.LoginRequest{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			mockBehavior: func(repo *MockUserRepository) {
				repo.On("FindByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Invalid Credentials",
			input: models.LoginRequest{
				Identifier: "test@example.com",
				Password:   "wrongpassword",
			},
			mockBehavior: func(repo *MockUserRepository) {
				repo.On("FindByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockBehavior(mockRepo)
			handler := NewUserHandler(mockRepo)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.Login(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

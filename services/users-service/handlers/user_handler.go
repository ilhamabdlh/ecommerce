package handlers

import (
	"net/http"
	"strings"

	"ecommerce/user-service/models"
	"ecommerce/user-service/repository"
	"ecommerce/user-service/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserHandler(userRepo repository.UserRepositoryInterface) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration details"
// @Success 201 {object} map[string]interface{} "message,user_id"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.LogError("Invalid registration request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError("Failed to hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Create user
	if err := h.userRepo.Create(c.Request.Context(), &user); err != nil {
		utils.LogError("Failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	utils.LogInfo("User registered successfully",
		zap.String("user_id", user.ID.Hex()),
		zap.String("email", user.Email))

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": user.ID})
}

// Login godoc
// @Summary Login user
// @Description Login with email/phone and password
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.TokenResponse
// @Failure 401 {object} map[string]interface{} "error"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var loginReq models.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari user berdasarkan email atau phone
	var user *models.User
	var err error

	if isEmail(loginReq.Identifier) {
		user, err = h.userRepo.FindByEmail(c.Request.Context(), loginReq.Identifier)
	} else {
		user, err = h.userRepo.FindByPhone(c.Request.Context(), loginReq.Identifier)
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func isEmail(identifier string) bool {
	// Simple email validation
	return strings.Contains(identifier, "@")
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the profile of the currently logged in user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]interface{} "error"
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile of the currently logged in user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body models.User true "User update details"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 401 {object} map[string]interface{} "error"
// @Router /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userRepo.Update(c.Request.Context(), userID.(string), &updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

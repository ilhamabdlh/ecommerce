package handlers

import (
	"ecommerce/warehouse-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	UserID string   `json:"user_id" binding:"required"`
	Roles  []string `json:"roles" binding:"required"`
}

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// @Summary Generate JWT token
// @Description Generate JWT token for testing
// @Tags auth
// @Accept json
// @Produce json
// @Param request body TokenRequest true "Token Request"
// @Success 200 {object} map[string]string
// @Router /api/v1/auth/token [post]
func (h *AuthHandler) GenerateToken(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(req.UserID, req.Roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

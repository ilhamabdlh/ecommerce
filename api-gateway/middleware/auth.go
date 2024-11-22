package middleware

import (
	"context"
	"net/http"

	pb "github.com/ilhamabdlh/ecommerce/proto"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userClient pb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		resp, err := userClient.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
			Token: token,
		})
		if err != nil || !resp.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", resp.UserId)
		c.Next()
	}
}

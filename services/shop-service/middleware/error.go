package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.Error("Request error",
					zap.String("error", e.Error()),
					zap.String("path", c.Request.URL.Path),
				)
			}

			// Return last error to client
			lastError := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": lastError.Error(),
			})
		}
	}
}

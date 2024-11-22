package middleware

import (
	"ecommerce/internal/pkg/resilience"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CircuitBreakerMiddleware(name string) gin.HandlerFunc {
	cb := resilience.NewCircuitBreaker(resilience.CircuitBreakerConfig{
		Name:        name,
		MaxRequests: 100,
		Interval:    time.Minute,
		Timeout:     time.Second * 30,
	})

	return func(c *gin.Context) {
		result, err := cb.Execute(func() (interface{}, error) {
			c.Next()
			if c.Writer.Status() >= 500 {
				return nil, fmt.Errorf("server error: %d", c.Writer.Status())
			}
			return nil, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "Service temporarily unavailable",
			})
			return
		}

		if result != nil {
			c.JSON(http.StatusOK, result)
		}
	}
}

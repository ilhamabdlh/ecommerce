package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	ips    map[string][]time.Time
	mu     sync.RWMutex
	limit  int
	window time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		ips:    make(map[string][]time.Time),
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		if _, exists := rl.ips[ip]; !exists {
			rl.ips[ip] = []time.Time{now}
			c.Next()
			return
		}

		requests := rl.ips[ip]
		windowStart := now.Add(-rl.window)

		// Remove old requests
		var valid []time.Time
		for _, t := range requests {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}

		if len(valid) >= rl.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		rl.ips[ip] = append(valid, now)
		c.Next()
	}
}

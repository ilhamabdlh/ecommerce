package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRateLimiter struct {
	client      *redis.Client
	maxRequests int
	duration    time.Duration
}

func NewRedisRateLimiter(client *redis.Client, maxRequests int, duration time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		client:      client,
		maxRequests: maxRequests,
		duration:    duration,
	}
}

func (rl *RedisRateLimiter) IsAllowed(ctx context.Context, key string) (bool, error) {
	pipe := rl.client.Pipeline()
	now := time.Now().UnixNano()
	key = fmt.Sprintf("ratelimit:%s", key)

	// Clean old requests
	clearBefore := time.Now().Add(-rl.duration).UnixNano()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprint(clearBefore))

	// Add current request
	pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})

	// Count requests
	pipe.ZCard(ctx, key)

	// Set expiry
	pipe.Expire(ctx, key, rl.duration)

	cmders, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	// Get the count from the third command (ZCard)
	count := cmders[2].(*redis.IntCmd).Val()
	return count <= int64(rl.maxRequests), nil
}

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

func (rl *RedisRateLimiter) Allow(key string) (bool, error) {
	ctx := context.Background()
	now := time.Now().Unix()
	key = fmt.Sprintf("ratelimit:%s", key)

	pipe := rl.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-int64(rl.duration.Seconds())))
	pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})
	pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, rl.duration)

	cmders, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count := cmders[2].(*redis.IntCmd).Val()
	return count <= int64(rl.maxRequests), nil
}

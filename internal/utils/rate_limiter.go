package utils

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	limits map[string]struct {
		requests int
		duration time.Duration
	}
}

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{
		client: client,
		limits: map[string]struct {
			requests int
			duration time.Duration
		}{
			"comment": {5, time.Minute},
			"vote":    {10, time.Minute},
		},
	}
}

func (r *RateLimiter) Check(ctx context.Context, userID, action string) error {
	limit, ok := r.limits[action]
	if !ok {
		return nil
	}

	key := "rate_limit:" + action + ":" + userID
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	if count == 1 {
		r.client.Expire(ctx, key, limit.duration)
	}

	if count > int64(limit.requests) {
		return errors.New("rate limit exceeded")
	}

	return nil
}

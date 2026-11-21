package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type rateLimitRepository struct {
	c *redis.Client
}

func NewRepository(c *redis.Client) Repository {
	return &rateLimitRepository{c: c}
}

func (r *rateLimitRepository) IncreaseRateLimit(ctx context.Context, key string, exp time.Duration) {
	r.c.Incr(ctx, key)
	r.c.ExpireNX(ctx, key, exp)
}
func (r *rateLimitRepository) GetCurrentRateLimit(ctx context.Context, key string) (int, error) {
	val, err := r.c.Get(ctx, key).Int()
	if err != nil {
		return -1, err
	}
	return val, nil
}

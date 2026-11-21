package ratelimit

import (
	"context"
	"time"
)

type Repository interface {
	IncreaseRateLimit(ctx context.Context, key string, exp time.Duration)
	GetCurrentRateLimit(ctx context.Context, key string) (int, error)
}

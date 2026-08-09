package cache

import (
	"context"
	"time"
)

//go:generate mockery --name=DB --filename=db.go --outpkg=mock_cache

type DB interface {
	SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error
	GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error)
	DeleteCacheKey(ctx context.Context, cacheGroupKey, cacheKey string) error
	DeleteCacheGroupKey(ctx context.Context, cacheGroupKey string) error
}

package url

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=URLRepository --filename=url.go

//--filename=../../mocks/healthcheck.go

// UrlStorage is the interface for URL storage
type URLRepository interface {
	StoreURL(ctx context.Context, code, url string, exp time.Duration) error
	GetURL(ctx context.Context, code string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

// NewUrlStorage creates a new UrlStorage
func NewRepository(c *redis.Client) URLRepository {
	return &urlStorage{c: c}
}

package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	urlExpTime = 24 * time.Hour
)

// UrlStorage is the interface for URL storage
type UrlStorage interface {
	StoreURL(ctx context.Context, code, url string) error
	GetURL(ctx context.Context, code string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

// NewUrlStorage creates a new UrlStorage
func NewUrlStorage(c *redis.Client) UrlStorage {
	return &urlStorage{c: c}
}

// StoreURL stores a URL in the cache
func (s *urlStorage) StoreURL(ctx context.Context, code, url string) error {
	return s.c.Set(ctx, code, url, urlExpTime).Err()
}

// GetURL retrieves a URL from the cache
func (s *urlStorage) GetURL(ctx context.Context, code string) (string, error) {
	return s.c.Get(ctx, code).Result()
}

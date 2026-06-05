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
	GetURL(ctx context.Context, code string) (string, error)
	StoreUrlIfUniqueCode(ctx context.Context, code string, url string, expireTime int) (bool, error)
}

type urlStorage struct {
	c *redis.Client
}

// NewUrlStorage creates a new UrlStorage
func NewUrlStorage(c *redis.Client) UrlStorage {
	return &urlStorage{c: c}
}

// GetURL retrieves a URL from the cache
func (s *urlStorage) GetURL(ctx context.Context, code string) (string, error) {
	return s.c.Get(ctx, code).Result()
}

// check for unique Code and save with expire time
func (s *urlStorage) StoreUrlIfUniqueCode(ctx context.Context, code string, url string, expireTime int) (bool, error) {
	timeDuration := time.Duration(expireTime) * time.Second
	success, err := s.c.SetNX(ctx, code, url, timeDuration).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

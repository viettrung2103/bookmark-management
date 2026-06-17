package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	urlExpTime = 24 * time.Hour
)

//go:generate mockery --name=UrlStorage --filename=urlstorage.go

// UrlStorage is the interface for URL storage
type UrlStorage interface {
	StoreURL(ctx context.Context, code, url string, exp time.Duration) error
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
func (s *urlStorage) StoreURL(ctx context.Context, code, url string, exp time.Duration) error {
	err := s.c.Set(ctx, code, url, exp).Err()
	if err != nil {
		log.Error().Err(err).Str("from", "repo.urlStorage.StoreURL").Msg("failed to store url")

		return err
	}
	return nil
}

// GetURL retrieves a URL from the cache
func (s *urlStorage) GetURL(ctx context.Context, code string) (string, error) {

	return s.c.Get(ctx, code).Result()
}

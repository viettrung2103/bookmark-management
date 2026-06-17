package urlstorage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

//go:generate mockery --name=UrlStorage --filename=urlstorage.go

// UrlStorage is the interface for URL storage
type Repository interface {
	StoreURL(ctx context.Context, code, url string, exp time.Duration) error
	GetURL(ctx context.Context, code string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

// NewUrlStorage creates a new UrlStorage
func NewRepository(c *redis.Client) Repository {
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

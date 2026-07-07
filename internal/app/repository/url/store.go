package url

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

// StoreURL stores a URL in the cache
func (s *urlStorage) StoreURL(ctx context.Context, code, url string, exp time.Duration) error {
	err := s.c.Set(ctx, code, url, exp).Err()
	if err != nil {
		log.Error().Err(err).Str("from", "repo.urlStorage.StoreURL").Msg("failed to store url")

		return err
	}
	return nil
}

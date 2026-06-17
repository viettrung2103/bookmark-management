package urlstorage

import (
	"context"
)

// GetURL retrieves a URL from the cache
func (s *urlStorage) GetURL(ctx context.Context, code string) (string, error) {

	return s.c.Get(ctx, code).Result()
}

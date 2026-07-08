package link

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var ErrCodeDoesNotExist = errors.New("code does not exist")

// GetLinkFromCode gets the url from the code
func (s *shortenUrlService) GetLinkFromCode(ctx context.Context, urlCode string) (string, error) {
	url, err := s.repo.GetURL(ctx, urlCode)
	if errors.Is(err, redis.Nil) {
		log.Error().Err(err).Str("from", "service.shortenUrlService.GetLinkFromCode").Msg("failed to get url from code")

		return "", ErrCodeDoesNotExist
	}

	return url, err
}

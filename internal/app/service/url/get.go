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
	println("get url from code", urlCode)
	if s.keygen.IsRedisCode(urlCode) {
		println("is redis code")
		url, err := s.urlRepo.GetURL(ctx, urlCode)
		if errors.Is(err, redis.Nil) {
			log.Error().Err(err).Str("from", "service.shortenUrlService.GetLinkFromCode").Msg("failed to get url from code")

			return "", ErrCodeDoesNotExist
		}
		return url, err
	}
	if s.keygen.IsDBCode(urlCode) {
		println("is db code")
		selectedBookmark, err := s.bookmarkRepo.GetBookmarkByCode(ctx, urlCode)
		if err != nil {
			return "", ErrCodeDoesNotExist
		}
		return selectedBookmark.URL, nil
	}
	return "", nil
}

package urlstorage

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	urlCodeLength = 7
)

// ShortenUrlWithExpiringTime shortens a url with expiring time
func (s *shortenUrlService) ShortenUrlWithExpiringTime(ctx context.Context, url string, expireTime int) (string, error) {
	// tao key
	key := s.keygen.GenerateKey(urlCodeLength)

	res, err := s.repo.GetURL(ctx, key)

	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	if res != "" {

		return s.ShortenUrlWithExpiringTime(ctx, url, expireTime)
	}

	// put key into redis
	err = s.repo.StoreURL(ctx, key, url, time.Duration(expireTime)*time.Second)
	if err != nil {
		log.Error().Err(err).Str("from", "service.shortenUrlService.ShortenUrlWithExpiringTime").Msg("failed to store url")

		return "", err
	}
	return key, nil
}

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

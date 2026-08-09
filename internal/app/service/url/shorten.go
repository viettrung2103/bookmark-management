package link

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	redisCodeLength = 7
)

// ShortenUrlWithExpiringTime shortens a url with expiring time
func (s *shortenUrlService) ShortenUrlWithExpiringTime(ctx context.Context, url string, expireTime int) (string, error) {
	// tao key
	key := s.keygen.GenerateRedisKey(redisCodeLength)

	res, err := s.urlRepo.GetURL(ctx, key)

	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	if res != "" {

		return s.ShortenUrlWithExpiringTime(ctx, url, expireTime)
	}

	// put key into redis
	err = s.urlRepo.StoreURL(ctx, key, url, time.Duration(expireTime)*time.Second)
	if err != nil {
		log.Error().Err(err).Str("from", "service.shortenUrlService.ShortenUrlWithExpiringTime").Msg("failed to store url")

		return "", err
	}
	return key, nil
}

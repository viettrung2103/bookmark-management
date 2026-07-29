package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisDB struct {
	c *redis.Client
}

func NewRedisDB(c *redis.Client) DB {
	return &redisDB{c: c}
}

func (r *redisDB) SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error {

	_, err := r.c.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, cacheGroupKey, cacheKey, value)
		pipe.Expire(ctx, cacheGroupKey, exp)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *redisDB) GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error) {
	result, err := r.c.HGet(ctx, cacheGroupKey, cacheKey).Bytes()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// delete all data of cache key
func (r *redisDB) DeleteCacheKey(ctx context.Context, cacheGroupKey, cacheKey string) error {
	err := r.c.HDel(ctx, cacheGroupKey, cacheKey).Err()
	if err != nil {
		return err
	}
	return nil
}

// delete all data of that group from top to bottom
func (r *redisDB) DeleteCacheGroupKey(ctx context.Context, cacheGroupKey string) error {
	err := r.c.Del(ctx, cacheGroupKey).Err()
	if err != nil {
		return err
	}
	return nil
}

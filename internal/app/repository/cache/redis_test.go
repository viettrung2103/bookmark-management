package cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/viettrung2103/bookmark-management/pkg/redis"
)

func TestRedisDB_SetCacheData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func() *redis.Client

		expectErr  error
		verifyFunc func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",
			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			expectErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				// Verify the data was stored in the Hash correctly
				val, err := r.HGet(ctx, "group_1", "key_1").Bytes()
				assert.Nil(t, err)
				assert.Equal(t, []byte("test_data"), val)

				// Optional: Verify the TTL/Expiration was set on the group
				ttl, _ := r.TTL(ctx, "group_1").Result()
				assert.Greater(t, ttl, time.Duration(0))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMock := tc.setupMock()
			testRepo := NewRedisDB(redisMock)

			err := testRepo.SetCacheData(ctx, "group_1", "key_1", []byte("test_data"), 10*time.Minute)
			assert.Equal(t, tc.expectErr, err)
			if err == nil {
				tc.verifyFunc(ctx, redisMock)
			}
		})
	}
}
func TestRedisDB_GetCacheData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *redis.Client

		expectResult []byte // Added to check the return value
		expectErr    error
	}{
		{
			name: "normal case - cache hit",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				// Seed the database with data before testing Get
				mock.HSet(ctx, "group_1", "key_1", []byte("test_data"))
				return mock
			},
			expectResult: []byte("test_data"),
			expectErr:    nil,
		},
		{
			name: "empty case - cache miss",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				// Return an empty db
				return mock
			},
			expectResult: nil,
			expectErr:    redis.Nil, // Expect Redis not found error
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMock := tc.setupMock(ctx)
			testRepo := NewRedisDB(redisMock)

			result, err := testRepo.GetCacheData(ctx, "group_1", "key_1")
			assert.Equal(t, tc.expectResult, result)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestRedisDB_DeleteCacheKey(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *redis.Client

		expectErr  error
		verifyFunc func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				// Seed it with multiple fields in the hash
				mock.HSet(ctx, "group_1", "key_1", "data_1")
				mock.HSet(ctx, "group_1", "key_2", "data_2")
				return mock
			},
			expectErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				// Verify key_1 is gone
				_, err := r.HGet(ctx, "group_1", "key_1").Result()
				assert.Equal(t, redis.Nil, err)

				// Verify key_2 is still there (since HDel only deletes specific fields)
				val, err := r.HGet(ctx, "group_1", "key_2").Result()
				assert.Nil(t, err)
				assert.Equal(t, "data_2", val)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMock := tc.setupMock(ctx)
			testRepo := NewRedisDB(redisMock)

			err := testRepo.DeleteCacheKey(ctx, "group_1", "key_1")
			assert.Equal(t, tc.expectErr, err)
			if err == nil {
				tc.verifyFunc(ctx, redisMock)
			}
		})
	}
}

func TestRedisDB_DeleteCacheGroupKey(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *redis.Client

		expectErr  error
		verifyFunc func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				// Seed it with multiple fields
				mock.HSet(ctx, "group_1", "key_1", "data_1")
				mock.HSet(ctx, "group_1", "key_2", "data_2")
				return mock
			},
			expectErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				// Verify the entire group/folder is gone
				exists, err := r.Exists(ctx, "group_1").Result()
				assert.Nil(t, err)
				assert.Equal(t, int64(0), exists) // 0 means it does not exist
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMock := tc.setupMock(ctx)
			testRepo := NewRedisDB(redisMock)

			err := testRepo.DeleteCacheGroupKey(ctx, "group_1")
			assert.Equal(t, tc.expectErr, err)
			if err == nil {
				tc.verifyFunc(ctx, redisMock)
			}
		})
	}
}

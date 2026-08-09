package url

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/viettrung2103/bookmark-management/pkg/redis"
)

// TestUrlStorage_StoreUrl tests the StoreUrl method
func TestUrlStorage_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		code      string
		setupMock func(ctx context.Context) *redis.Client
		//setupDB   func(ctx context.Context) *gorm.DB

		expectedVal string
		expectedErr error
	}{
		{
			name: "normal case",

			code: "1234567",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)

				// Pre-seed the cache with our test data so GetURL can find it
				err := mock.Set(ctx, "1234567", "https://google.com", 10*time.Minute).Err()
				if err != nil {
					t.Fatalf("failed to seed mock redis: %v", err)
				}
				return mock
			},
			expectedVal: "https://google.com",
			expectedErr: nil,
		},
		{
			name: "error - code does not exist (cache miss)",
			code: "nonexistent_code",
			setupMock: func(ctx context.Context) *redis.Client {
				// Return an empty mock database
				return redisPkg.InitMockRedis(t)
			},
			expectedVal: "",
			// go-redis returns redis.Nil when a key does not exist
			expectedErr: redis.Nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMock := tc.setupMock(ctx)

			testRepo := NewRepository(redisMock)

			val, err := testRepo.GetURL(ctx, tc.code)

			assert.Equal(t, tc.expectedVal, val)

			if tc.expectedErr != nil {
				assert.True(t, errors.Is(err, tc.expectedErr))
			} else {
				assert.NoError(t, err)
			}
		})

	}
}

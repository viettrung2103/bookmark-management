package repository

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redigPkg "github.com/viettrung2103/bookmark-management/pkg/redis"
)

func TestUrlStorage_StoreUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		expireDuration int

		setupMock func() *redis.Client

		verifyFunc   func(ctx context.Context, r *redis.Client)
		expectErr    error
		expectedBool bool
	}{
		{
			name:           "normal case",
			expireDuration: 1000,

			setupMock: func() *redis.Client {
				mock := redigPkg.InitMockRedis(t)
				return mock
			},
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				url, err := r.Get(ctx, "1234567").Result()
				assert.Nil(t, err)
				assert.Equal(t, url, "https://google.com")
			},
			expectErr:    nil,
			expectedBool: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMock := tc.setupMock()
			testRepo := NewUrlStorage(redisMock)

			bool, err := testRepo.StoreUrlIfUniqueCode(ctx, "1234567", "https://google.com", tc.expireDuration)
			assert.Equal(t, tc.expectErr, err)
			if err == nil {
				tc.verifyFunc(ctx, redisMock)
			}
			assert.Equal(t, tc.expectedBool, bool, tc.expectedBool)
		})

	}
}

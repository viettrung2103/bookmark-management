package service

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/repository"
	redigPkg "github.com/viettrung2103/bookmark-management/pkg/redis"
)

const (
	testUrl        = "https://google.com"
	expireDuration = 1000
)

func TestShortenUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		setupMock func() *redis.Client

		expectedUrl string
	}{
		{
			name: "successfull case",
			setupMock: func() *redis.Client {
				mock := redigPkg.InitMockRedis(t)
				return mock
			},
			expectedUrl: testUrl,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testMock := tc.setupMock()
			testRepo := repository.NewUrlStorage(testMock)
			testSvc := NewShortenUrl(testRepo)
			result, err := testSvc.ShortenUrlWithExpiringTime(context.Background(), testUrl, expireDuration)
			returnUrl, err := testRepo.GetURL(context.Background(), result)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUrl, returnUrl)

		})
	}
}

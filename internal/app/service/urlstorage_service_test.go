package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	repoMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/mocks"

	keygenMock "github.com/viettrung2103/bookmark-management/pkg/stringutils/mocks"
)

var redisTestErr = errors.New("test error")

func TestService_GetLinkFromKey(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRepo   func(ctx context.Context) *repoMocks.UrlStorage
		expectedUrl string
		expectedErr error
	}{
		{
			name: "normal case",
			setupRepo: func(ctx context.Context) *repoMocks.UrlStorage {
				mock := repoMocks.NewUrlStorage(t)
				mock.On("GetURL", ctx, "test").Return("https://test.com", nil)
				return mock
			},

			expectedUrl: "https://test.com",
			expectedErr: nil,
		},
		{
			name: "empty case",
			setupRepo: func(ctx context.Context) *repoMocks.UrlStorage {
				mock := repoMocks.NewUrlStorage(t)
				mock.On("GetURL", ctx, "test").Return("", redisTestErr)
				return mock
			},
			expectedUrl: "",
			expectedErr: redisTestErr,
		},
		{
			name: "err case ",
			setupRepo: func(ctx context.Context) *repoMocks.UrlStorage {
				mock := repoMocks.NewUrlStorage(t)
				mock.On("GetURL", ctx, "test").Return("", redisTestErr)
				return mock
			},

			expectedUrl: "",
			expectedErr: redisTestErr,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			mockRepo := tc.setupRepo(ctx)
			testService := NewShortenUrl(mockRepo, nil)
			result, err := testService.GetLinkFromCode(ctx, "test")
			assert.Equal(t, result, tc.expectedUrl)
			assert.ErrorIs(t, err, tc.expectedErr)

		})
	}
}

const testExpTime = 60 * time.Second

const linkKeyLength = 7

func TestService_CreateShortenLink(t *testing.T) {
	testCases := []struct {
		name        string
		setupRepo   func(ctx context.Context) *repoMocks.UrlStorage
		setupKeyGen func() *keygenMock.KeyGenerator

		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case - new key",

			setupRepo: func(ctx context.Context) *repoMocks.UrlStorage {
				mock := repoMocks.NewUrlStorage(t)
				mock.On("GetURL", ctx, "1234567").Return("", redis.Nil)
				mock.On("StoreURL", ctx, "1234567", "https://test.com", testExpTime).Return(nil)

				return mock

			},
			setupKeyGen: func() *keygenMock.KeyGenerator {
				mockKeyGen := keygenMock.NewKeyGenerator(t)
				mockKeyGen.On("GenerateKey", linkKeyLength).Return("1234567")

				return mockKeyGen
			},

			expectedErr:    nil,
			expectedResult: "1234567",
		},
		{
			name: "normal case - random the same key",

			setupRepo: func(ctx context.Context) *repoMocks.UrlStorage {
				mock := repoMocks.NewUrlStorage(t)
				// generate a key >> return a url >> generate new key
				mock.On("GetURL", ctx, "1234567").Return("https://example.com", redis.Nil)
				mock.On("GetURL", ctx, "2345678").Return("", redis.Nil)
				// store new key
				mock.On("StoreURL", ctx, "2345678", "https://test.com", testExpTime).Return(nil)

				return mock

			},
			setupKeyGen: func() *keygenMock.KeyGenerator {
				mockKeyGen := keygenMock.NewKeyGenerator(t)
				mockKeyGen.On("GenerateKey", linkKeyLength).Return("1234567").Once()
				mockKeyGen.On("GenerateKey", linkKeyLength).Return("2345678").Once()

				return mockKeyGen
			},

			expectedErr:    nil,
			expectedResult: "2345678",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			mockRepo := tc.setupRepo(ctx)
			keygenMock := tc.setupKeyGen()

			testService := NewShortenUrl(mockRepo, keygenMock)
			result, err := testService.ShortenUrlWithExpiringTime(ctx, "https://test.com", 60)
			assert.Equal(t, result, tc.expectedResult)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

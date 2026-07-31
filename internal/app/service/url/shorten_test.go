package link

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	bookmarkMockRepo "github.com/viettrung2103/bookmark-management/internal/app/repository/bookmark/mocks"
	urlMockRepo "github.com/viettrung2103/bookmark-management/internal/app/repository/url/mocks"
	//keygenMock"github.com/viettrung2103/bookmark-management/internal/app/service/bookmark/mocks"

	keygenMock "github.com/viettrung2103/bookmark-management/pkg/stringutils/mocks"
)

const testExpTime = 60 * time.Second

const linkKeyLength = 7

// TestService_CreateShortenLink tests the CreateShortenLink method
func TestService_CreateShortenLink(t *testing.T) {
	testCases := []struct {
		name              string
		setupURLRepo      func(ctx context.Context) *urlMockRepo.URLRepository
		setupBookmarkRepo func(ctx context.Context) *bookmarkMockRepo.Repository
		setupKeyGen       func() *keygenMock.KeyGenerator

		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case - new key",

			setupURLRepo: func(ctx context.Context) *urlMockRepo.URLRepository {
				mock := urlMockRepo.NewURLRepository(t)
				mock.On("GetURL", ctx, "1234567").Return("", redis.Nil)
				mock.On("StoreURL", ctx, "1234567", "https://test.com", testExpTime).Return(nil)

				return mock

			},
			setupBookmarkRepo: func(ctx context.Context) *bookmarkMockRepo.Repository {
				mock := bookmarkMockRepo.NewRepository(t)
				return mock
			},
			setupKeyGen: func() *keygenMock.KeyGenerator {
				mockKeyGen := keygenMock.NewKeyGenerator(t)
				mockKeyGen.On("GenerateRedisKey", linkKeyLength).Return("1234567")
				return mockKeyGen
			},

			expectedErr:    nil,
			expectedResult: "1234567",
		},
		{
			name: "normal case - random the same key",

			setupURLRepo: func(ctx context.Context) *urlMockRepo.URLRepository {
				mock := urlMockRepo.NewURLRepository(t)
				// generate a key >> return a url >> generate new key
				mock.On("GetURL", ctx, "1234567").Return("https://example.com", redis.Nil)
				mock.On("GetURL", ctx, "2345678").Return("", redis.Nil)
				// store new key
				mock.On("StoreURL", ctx, "2345678", "https://test.com", testExpTime).Return(nil)

				return mock

			},
			setupBookmarkRepo: func(ctx context.Context) *bookmarkMockRepo.Repository {
				mock := bookmarkMockRepo.NewRepository(t)
				return mock
			},
			setupKeyGen: func() *keygenMock.KeyGenerator {
				mockKeyGen := keygenMock.NewKeyGenerator(t)
				mockKeyGen.On("GenerateRedisKey", linkKeyLength).Return("1234567").Once()
				mockKeyGen.On("GenerateRedisKey", linkKeyLength).Return("2345678").Once()

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
			ctx := t.Context()

			urlRepoMock := tc.setupURLRepo(ctx)
			bookmarkRepoMock := tc.setupBookmarkRepo(ctx)

			keygen := tc.setupKeyGen()

			testService := NewService(urlRepoMock, bookmarkRepoMock, keygen)
			result, err := testService.ShortenUrlWithExpiringTime(ctx, "https://test.com", 60)
			assert.Equal(t, result, tc.expectedResult)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

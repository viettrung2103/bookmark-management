package link

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	bookmarkMockRepo "github.com/viettrung2103/bookmark-management/internal/app/repository/bookmark/mocks"
	keygenMocks "github.com/viettrung2103/bookmark-management/pkg/stringutils/mocks"

	URLMockRepo "github.com/viettrung2103/bookmark-management/internal/app/repository/url/mocks"
)

var redisTestErr = errors.New("test error")

// TestService_GetLinkFromKey tests the GetLinkFromKey method
func TestService_GetLinkFromKey(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupURLRepo      func(ctx context.Context) *URLMockRepo.URLRepository
		setupBookmarkRepo func(ctx context.Context) *bookmarkMockRepo.Repository
		setupKeygen       func(ctx context.Context) *keygenMocks.KeyGenerator
		expectedUrl       string
		expectedErr       error
	}{
		{
			name: "normal case",
			setupURLRepo: func(ctx context.Context) *URLMockRepo.URLRepository {
				mock := URLMockRepo.NewURLRepository(t)
				mock.On("GetURL", ctx, "test").Return("https://test.com", nil)
				return mock
			},
			setupBookmarkRepo: func(ctx context.Context) *bookmarkMockRepo.Repository {
				mock := bookmarkMockRepo.NewRepository(t)
				//mock.On("GetURL", ctx, "test").Return("https://test.com", nil)
				return mock
			},
			setupKeygen: func(ctx context.Context) *keygenMocks.KeyGenerator {
				mock := keygenMocks.NewKeyGenerator(t)
				mock.On("IsRedisCode", "test").Return(true)
				return mock
			},

			expectedUrl: "https://test.com",
			expectedErr: nil,
		},
		{
			name: "empty case",
			setupURLRepo: func(ctx context.Context) *URLMockRepo.URLRepository {
				mock := URLMockRepo.NewURLRepository(t)
				mock.On("GetURL", ctx, "test").Return("", redisTestErr)
				return mock
			},
			setupBookmarkRepo: func(ctx context.Context) *bookmarkMockRepo.Repository {
				mock := bookmarkMockRepo.NewRepository(t)
				//mock.On("GetURL", ctx, "test").Return("https://test.com", nil)
				return mock
			},
			setupKeygen: func(ctx context.Context) *keygenMocks.KeyGenerator {
				mock := keygenMocks.NewKeyGenerator(t)
				mock.On("IsRedisCode", "test").Return(true)
				return mock
			},
			expectedUrl: "",
			expectedErr: redisTestErr,
		},
		{
			name: "err case ",
			setupURLRepo: func(ctx context.Context) *URLMockRepo.URLRepository {
				mock := URLMockRepo.NewURLRepository(t)
				mock.On("GetURL", ctx, "test").Return("", redisTestErr)
				return mock
			},
			setupBookmarkRepo: func(ctx context.Context) *bookmarkMockRepo.Repository {
				mock := bookmarkMockRepo.NewRepository(t)
				//mock.On("GetURL", ctx, "test").Return("https://test.com", nil)
				return mock
			},
			setupKeygen: func(ctx context.Context) *keygenMocks.KeyGenerator {
				mock := keygenMocks.NewKeyGenerator(t)
				mock.On("IsRedisCode", "test").Return(true)
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
			ctx := t.Context()

			urlMock := tc.setupURLRepo(ctx)
			keygenMock := tc.setupKeygen(ctx)
			bookmarkMock := tc.setupBookmarkRepo(ctx)
			testService := NewService(urlMock, bookmarkMock, keygenMock)
			result, err := testService.GetLinkFromCode(ctx, "test")
			assert.Equal(t, result, tc.expectedUrl)
			assert.ErrorIs(t, err, tc.expectedErr)

		})
	}
}

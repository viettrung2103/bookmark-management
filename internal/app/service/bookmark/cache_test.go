package bookmark_test

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	mock_cache "github.com/viettrung2103/bookmark-management/internal/app/repository/cache/mocks"
	"github.com/viettrung2103/bookmark-management/internal/app/service/bookmark"
	"github.com/viettrung2103/bookmark-management/internal/app/service/bookmark/mocks"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
)

const (
	validCacheGroupKey = "get_bookmarks_133b3b42-70b9-456c-82e7-bf1b570e6c51"
	validUserID        = "133b3b42-70b9-456c-82e7-bf1b570e6c51"
	validCacheKey      = "page_1_limit_10"
)

func TestBookmarkCacheService_GetBookmarks(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name string

		setupService func(ctx context.Context) *mocks.Service
		setupCache   func(ctx context.Context) *mock_cache.DB

		expectedResult *bookmark.GetBookmarkResult
		expectedError  error
	}{
		{
			name: "success - with cache",
			setupService: func(ctx context.Context) *mocks.Service {
				mockService := mocks.NewService(t)
				return mockService
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockCache := mock_cache.NewDB(t)
				mockCache.On("GetCacheData", ctx, validCacheGroupKey, validCacheKey).Return(
					[]byte(`{"bookmarks":[{"id":"f4defa89-c5b3-4f26-8ca1-614495cdde12","description":"Google","url":"https://www.google.com","code":"kAQllBal"},{"id":"f4defa89-c5b3-4f26-8ca1-614495cdde22","description":"Google1","url":"https://www.google1.com","code":"kAQllBal1"}],"count":2}`),
					nil,
				)
				return mockCache
			},
			expectedResult: &bookmark.GetBookmarkResult{
				Bookmarks: []*model.Bookmark{
					{
						Base:        fixtures.GetTestBase("f4defa89-c5b3-4f26-8ca1-614495cdde12"),
						Description: "Google",
						URL:         "https://www.google.com",
						Code:        "kAQllBal",
					},
					{
						Base:        fixtures.GetTestBase("f4defa89-c5b3-4f26-8ca1-614495cdde22"),
						Description: "Google1",
						URL:         "https://www.google1.com",
						Code:        "kAQllBal1",
					},
				},
				Count: 2,
			},
		},
		{
			name: "success - without cache",
			setupService: func(ctx context.Context) *mocks.Service {
				mockService := mocks.NewService(t)
				mockService.On("GetBookmarks", ctx, "133b3b42-70b9-456c-82e7-bf1b570e6c51", 1, 10).Return(
					&bookmark.GetBookmarkResult{
						Bookmarks: []*model.Bookmark{
							{
								Base:        fixtures.GetTestBase("f4defa89-c5b3-4f26-8ca1-614495cdde12"),
								Description: "Google",
								URL:         "https://www.google.com",
								Code:        "kAQllBal",
							},
							{
								Base:        fixtures.GetTestBase("f4defa89-c5b3-4f26-8ca1-614495cdde22"),
								Description: "Google1",
								URL:         "https://www.google1.com",
								Code:        "kAQllBal1",
							},
						},
						Count: 2,
					},
					nil)
				return mockService
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockCache := mock_cache.NewDB(t)
				mockCache.On("GetCacheData", ctx, validCacheGroupKey, validCacheKey).Return(
					nil,
					redis.Nil,
				)
				mockCache.On("SetCacheData", ctx, validCacheGroupKey, validCacheKey,
					[]byte(`{"bookmarks":[{"id":"f4defa89-c5b3-4f26-8ca1-614495cdde12","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","description":"Google","url":"https://www.google.com","code":"kAQllBal"},{"id":"f4defa89-c5b3-4f26-8ca1-614495cdde22","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","description":"Google1","url":"https://www.google1.com","code":"kAQllBal1"}],"count":2}`),
					bookmark.CacheExpireDuration,
				).Return(
					nil,
				)
				return mockCache
			},
			expectedResult: &bookmark.GetBookmarkResult{
				Bookmarks: []*model.Bookmark{
					{
						Base:        fixtures.GetTestBase("f4defa89-c5b3-4f26-8ca1-614495cdde12"),
						Description: "Google",
						URL:         "https://www.google.com",
						Code:        "kAQllBal",
					},
					{
						Base:        fixtures.GetTestBase("f4defa89-c5b3-4f26-8ca1-614495cdde22"),
						Description: "Google1",
						URL:         "https://www.google1.com",
						Code:        "kAQllBal1",
					},
				},
				Count: 2,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			mockService := tc.setupService(ctx)
			mockCache := tc.setupCache(ctx)

			bookmarkCacheService := bookmark.NewBookmarkCacheService(mockService, mockCache)
			res, err := bookmarkCacheService.GetBookmarks(ctx, validUserID, 1, 10)
			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

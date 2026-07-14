package bookmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	"gorm.io/gorm"
)

const (
	fixtureGetUserID      = "133b3b42-70b9-456c-82e7-bf1b570e6c51" // Owns exactly 2 bookmarks
	unauthorizedGetUserID = "133b3b42-70b9-456c-82e7-bf1b570e6c60" // Valid user, but owns 0 bookmarks
)

func TestBookmarkRepository_GetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB
		userID  string
		offset  int
		limit   int

		expectedCount int
		expectedURLs  []string
		expectedError error
	}{
		{
			name: "Get all bookmarks successfully (No pagination truncation)",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			userID:        fixtureGetUserID,
			offset:        0,
			limit:         10,
			expectedCount: 2,
			expectedURLs:  []string{"http://google.com", "http://google.com"},
			expectedError: nil,
		},
		{
			name: "Get bookmarks with pagination (Limit to 1)",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			userID:        fixtureGetUserID,
			offset:        0,
			limit:         1,
			expectedCount: 1,
			expectedError: nil,
		},
		{
			name: "Get bookmarks offset check (Skip first bookmark)",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			userID:        fixtureGetUserID,
			offset:        1,
			limit:         1,
			expectedCount: 1,
			expectedError: nil,
		},
		{
			name: "Get bookmarks for user with no data",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			userID:        unauthorizedGetUserID,
			offset:        0,
			limit:         10,
			expectedCount: 0,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewRepository(db)

			bookmarks, err := repo.GetBookmarks(ctx, tc.userID, tc.offset, tc.limit)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, bookmarks, tc.expectedCount)

			// If specific urls are specified, check ordering matches (created_at ASC)
			if len(tc.expectedURLs) > 0 {
				for i, expectedURL := range tc.expectedURLs {
					assert.Equal(t, expectedURL, bookmarks[i].URL)
				}
			}
		})
	}
}

func TestBookmarkRepository_GetBookmarkCount(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB
		userID  string

		expectedCount int64
		expectedError error
	}{
		{
			name: "Get correct count for active user with bookmarks",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			userID:        fixtureGetUserID,
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "Get zero count for user with no bookmarks",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			userID:        unauthorizedGetUserID,
			expectedCount: 0,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewRepository(db)

			count, err := repo.GetBookmarkCount(ctx, tc.userID)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCount, count)
		})
	}
}

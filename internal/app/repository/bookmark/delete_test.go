package bookmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

// Match these EXACTLY to your BookmarkCommonTestDB generated data records
const (
	fixtureUserID          = "133b3b42-70b9-456c-82e7-bf1b570e6c51" // Owns bookmark1 and bookmark2
	fixtureBookmarkID1     = "133b3b42-70b9-456c-93d8-bf1b570e6c55" // bookmark1
	nonExistingBookmarkID  = "133b3b42-70b9-456c-93d8-bf1b570e6c99" // Valid format, missing from DB
	unauthorizedTestUserID = "133b3b42-70b9-456c-82e7-bf1b570e6c60" // Valid user format, but doesn't own the bookmark
)

func TestBookmarkRepository_DeleteBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB    func(t *testing.T) *gorm.DB
		userID     string
		bookmarkID string

		expectedError error
	}{
		{
			name: "Delete bookmark successfully",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			userID:        fixtureUserID,
			bookmarkID:    fixtureBookmarkID1,
			expectedError: nil,
		},
		{
			name: "Error - Bookmark not found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			userID:        fixtureUserID,
			bookmarkID:    nonExistingBookmarkID,
			expectedError: dbutils.ErrRecordNotFound,
		},
		{
			name: "Error - Wrong user tries to delete another user's bookmark",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			userID:        unauthorizedTestUserID,
			bookmarkID:    fixtureBookmarkID1,        // Bookmark exists, but belongs to fixtureUserID
			expectedError: dbutils.ErrRecordNotFound, // RowsAffected will be 0 -> triggers NotFound
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewRepository(db)

			err := repo.DeleteBookmarkByID(ctx, tc.userID, tc.bookmarkID)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)

			// Double check mutation validation: Verify from the DB instance that it's completely missing
			//var count int64
			//db.WithContext(ctx).
			//	Model(&model.Bookmark{}).
			//	Where("id = ? AND user_id = ?", tc.bookmarkID, tc.userID).
			//	Count(&count)
			//
			//assert.Equal(t, int64(0), count, "Expected bookmark record to be completely missing from the database context")
		})
	}
}

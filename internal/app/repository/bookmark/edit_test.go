package bookmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils/mocks"
	"gorm.io/gorm"
)

const (
	fixtureEditUserID         = "133b3b42-70b9-456c-82e7-bf1b570e6c51" // Owns bookmark1 and bookmark2
	fixtureEditBookmarkID1    = "133b3b42-70b9-456c-93d8-bf1b570e6c55" // bookmark1
	nonExistingEditBookmarkID = "133b3b42-70b9-456c-93d8-bf1b570e6c99" // Missing from DB
	unauthorizedEditUserID    = "133b3b42-70b9-456c-82e7-bf1b570e6c60" // Does not own fixtureEditBookmarkID1
)

func TestBookmarkRepository_EditBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB        func(t *testing.T) *gorm.DB
		userID         string
		bookmarkID     string
		newDescription string
		newURL         string

		expectedError error
	}{
		{
			name: "Edit bookmark successfully",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			userID:         fixtureEditUserID,
			bookmarkID:     fixtureEditBookmarkID1,
			newDescription: "Fully Updated Description",
			newURL:         "https://updated-link.com",
			expectedError:  nil,
		},
		{
			name: "Error - Bookmark not found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			userID:         fixtureEditUserID,
			bookmarkID:     nonExistingEditBookmarkID,
			newDescription: "Some Description",
			newURL:         "https://some-url.com",
			expectedError:  dbutils.ErrRecordNotFound,
		},
		{
			name: "Error - Wrong user tries to edit another user's bookmark",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			userID:         unauthorizedEditUserID,
			bookmarkID:     fixtureEditBookmarkID1, // Belongs to fixtureEditUserID, not unauthorizedEditUserID
			newDescription: "Malicious Edit Attempt",
			newURL:         "https://hacked.com",
			expectedError:  dbutils.ErrRecordNotFound, // RowsAffected is 0 -> triggers NotFound
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			mockKeygen := mocks.NewKeyGenerator(t)
			repo := NewRepository(db, mockKeygen)

			err := repo.EditBookmarkByID(ctx, tc.userID, tc.bookmarkID, tc.newDescription, tc.newURL)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)

			// State Verification Check: Retrieve the record directly from GORM to assert updates occurred
			var updatedRecord model.Bookmark
			err = db.WithContext(ctx).
				Model(&model.Bookmark{}).
				Where("id = ? AND user_id = ?", tc.bookmarkID, tc.userID).
				First(&updatedRecord).
				Error

			assert.NoError(t, err, "Expected to find the updated bookmark in the database")
			assert.Equal(t, tc.newDescription, updatedRecord.Description, "Database description must match the updated value")
			assert.Equal(t, tc.newURL, updatedRecord.URL, "Database URL must match the updated value")
		})
	}
}

package bookmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

const existingTestUserID = "133b3b42-70b9-456c-82e7-bf1b570e6c52"
const nonExistingTestUserId = "133b3b42-70b9-456c-82e7-bf1b570e6c60"

func TestBookmarkRepository_CreateBookmark(t *testing.T) {
	t.Parallel()
	//testUserID :="87a3cb94-d2e8-422d-bb91-fc5215949eb8"),

	testCases := []struct {
		name string

		setupDB       func(t *testing.T) *gorm.DB
		inputBookmark *model.Bookmark

		expectedOutput *model.Bookmark
		expectedError  error
	}{
		{
			name: "Create bookmark successfully",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			inputBookmark: &model.Bookmark{
				Base:        fixtures.GetTestBase("133b3b42-70b9-456c-82e7-bf1b570e6c57"),
				UserID:      fixtures.GetUUID(existingTestUserID),
				URL:         "https://example/newbookmark",
				Description: "New bookmark for test user 1",
			},

			expectedOutput: &model.Bookmark{
				Base:        fixtures.GetTestBase("133b3b42-70b9-456c-82e7-bf1b570e6c57"),
				UserID:      fixtures.GetUUID(existingTestUserID),
				URL:         "https://example/newbookmark",
				Description: "New bookmark for test user 1",
			},
			//expectedError: nil,
		},
		{
			name: "No user",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkPlusTestDB{})
			},
			inputBookmark: &model.Bookmark{
				Base:        fixtures.GetTestBase("133b3b42-70b9-456c-82e7-bf1b570e6c60"),
				UserID:      fixtures.GetUUID(nonExistingTestUserId),
				URL:         "https://example/newbookmark",
				Description: "New bookmark for test user 1",
			},

			//expectedOutput: &model.Bookmark{},
			expectedError: dbutils.ErrForeignKey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewRepository(db)

			output, err := repo.CreateBookmark(ctx, tc.inputBookmark)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			assert.NoError(t, err)

			if err != nil {
				assert.Equal(t, tc.expectedOutput, output)
				return
			}

			assert.Equal(t, tc.expectedOutput.URL, output.URL)
			assert.Equal(t, tc.expectedOutput.Description, output.Description)
			assert.Equal(t, tc.expectedOutput.Code, output.Code)
		})

	}
}

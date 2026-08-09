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

const existingTestUserID = "133b3b42-70b9-456c-82e7-bf1b570e6c52"
const nonExistingTestUserId = "133b3b42-70b9-456c-82e7-bf1b570e6c60"
const mockCodeInt uint64 = 3

func TestBookmarkRepository_CreateBookmark(t *testing.T) {
	t.Parallel()
	//testUserID :="87a3cb94-d2e8-422d-bb91-fc5215949eb8"),

	testCases := []struct {
		name string

		setupDB       func(t *testing.T) *gorm.DB
		setupMock     func(keygenMock *mocks.KeyGenerator)
		inputBookmark *model.Bookmark

		expectedOutput *model.Bookmark
		expectedError  error
	}{
		{
			name: "Create bookmark successfully",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupMock: func(keygenMock *mocks.KeyGenerator) {
				// Import "github.com/stretchr/testify/mock" at the top of your file
				// Tell the mock to return "mockedCode" whenever GenerateBase62Code is called
				keygenMock.On("GenerateBase62Code", mockCodeInt).Return("0x3")
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
				Code:        "0x3",
			},
			//expectedError: nil,
		},
		{
			name: "No user",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkPlusTestDB{})
			},
			setupMock: func(keygenMock *mocks.KeyGenerator) {
				// Do nothing. The DB insert will fail first,
				// so GenerateBase62Code will never be called.
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
			mocKeygen := mocks.NewKeyGenerator(t)

			if tc.setupMock != nil {
				tc.setupMock(mocKeygen)
			}

			repo := NewRepository(db, mocKeygen)

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

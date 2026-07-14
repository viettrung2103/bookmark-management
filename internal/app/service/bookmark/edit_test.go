package bookmark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	repoMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/bookmark/mocks"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
)

const (
	testSvcEditUserID     = "133b3b42-70b9-456c-82e7-bf1b570e6c51"
	testSvcEditBookmarkID = "133b3b42-70b9-456c-93d8-bf1b570e6c55"
)

func TestBookmarkService_EditBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		userID         string
		bookmarkID     string
		newDescription string
		newURL         string
		setupMocks     func(ctx context.Context, repo *repoMocks.Repository)
		expectedError  error
	}{
		{
			name:           "Success - Edit Bookmark",
			userID:         testSvcEditUserID,
			bookmarkID:     testSvcEditBookmarkID,
			newDescription: "Updated Google Link",
			newURL:         "https://updated-google.com",
			setupMocks: func(ctx context.Context, repo *repoMocks.Repository) {
				// Verify service correctly forwards arguments down to repository layer
				repo.On("EditBookmarkByID", ctx, testSvcEditUserID, testSvcEditBookmarkID, "Updated Google Link", "https://updated-google.com").
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:           "Error - Bookmark Not Found",
			userID:         testSvcEditUserID,
			bookmarkID:     testSvcEditBookmarkID,
			newDescription: "Updated Google Link",
			newURL:         "https://updated-google.com",
			setupMocks: func(ctx context.Context, repo *repoMocks.Repository) {
				repo.On("EditBookmarkByID", ctx, testSvcEditUserID, testSvcEditBookmarkID, "Updated Google Link", "https://updated-google.com").
					Return(dbutils.ErrRecordNotFound)
			},
			expectedError: dbutils.ErrRecordNotFound,
		},
		{
			name:           "Error - Repository Duplication Error",
			userID:         testSvcEditUserID,
			bookmarkID:     testSvcEditBookmarkID,
			newDescription: "Duplicate Link",
			newURL:         "https://duplicate-url.com",
			setupMocks: func(ctx context.Context, repo *repoMocks.Repository) {
				repo.On("EditBookmarkByID", ctx, testSvcEditUserID, testSvcEditBookmarkID, "Duplicate Link", "https://duplicate-url.com").
					Return(dbutils.ErrDuplication)
			},
			expectedError: dbutils.ErrDuplication,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			// Initialize mock controllers
			mockRepo := repoMocks.NewRepository(t)
			tc.setupMocks(ctx, mockRepo)

			// Construct service injecting the mocked repository
			service := &bookmarkService{
				bookmarkRepo: mockRepo,
			}

			err := service.EditBookmarkByID(ctx, tc.userID, tc.bookmarkID, tc.newDescription, tc.newURL)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)
		})
	}
}

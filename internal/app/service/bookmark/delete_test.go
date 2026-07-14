package bookmark

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	RepoMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/bookmark/mocks"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
)

const (
	testSvcDeleteUserID     = "133b3b42-70b9-456c-82e7-bf1b570e6c51"
	testSvcDeleteBookmarkID = "133b3b42-70b9-456c-93d8-bf1b570e6c55"
)

func TestBookmarkService_DeleteBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		userID        string
		bookmarkID    string
		setupMocks    func(repo *RepoMocks.Repository)
		expectedError error
	}{
		{
			name:       "Success - Delete Bookmark",
			userID:     testSvcDeleteUserID,
			bookmarkID: testSvcDeleteBookmarkID,
			setupMocks: func(repo *RepoMocks.Repository) {
				// Expecting service layer to pass context, userID string, and bookmarkID string down directly
				repo.On("DeleteBookmarkByID", mock.Anything, testSvcDeleteUserID, testSvcDeleteBookmarkID).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "Error - Bookmark Not Found",
			userID:     testSvcDeleteUserID,
			bookmarkID: testSvcDeleteBookmarkID,
			setupMocks: func(repo *RepoMocks.Repository) {
				repo.On("DeleteBookmarkByID", mock.Anything, testSvcDeleteUserID, testSvcDeleteBookmarkID).
					Return(dbutils.ErrRecordNotFound)
			},
			expectedError: dbutils.ErrRecordNotFound,
		},
		{
			name:       "Error - Internal Database Failure",
			userID:     testSvcDeleteUserID,
			bookmarkID: testSvcDeleteBookmarkID,
			setupMocks: func(repo *RepoMocks.Repository) {
				repo.On("DeleteBookmarkByID", mock.Anything, testSvcDeleteUserID, testSvcDeleteBookmarkID).
					Return(errors.New("unexpected table deadlocks"))
			},
			expectedError: errors.New("unexpected table deadlocks"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Initialize mock controllers
			mockRepo := RepoMocks.NewRepository(t)
			tc.setupMocks(mockRepo)

			// Construct service injecting the mocked repository layer
			service := &bookmarkService{
				bookmarkRepo: mockRepo,
			}

			err := service.DeleteBookmarkByID(context.Background(), tc.userID, tc.bookmarkID)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

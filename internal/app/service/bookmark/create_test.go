package bookmark

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	repoMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/bookmark/mocks"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	keygenMocks "github.com/viettrung2103/bookmark-management/pkg/stringutils/mocks"
)

const (
	testSvcUserID     = "133b3b42-70b9-456c-82e7-bf1b570e6c51"
	mockGeneratedCode = "abc12345"
	mockBase62Code    = "base62xyz"
)

func TestBookmarkService_AddBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		description string
		url         string
		userID      string
		setupMocks  func(ctx context.Context, repo *repoMocks.Repository, keygen *keygenMocks.KeyGenerator)

		expectedOut *model.Bookmark
		expectedErr error
	}{
		{
			name:        "Success - Add Bookmark",
			description: "Google Search",
			url:         "https://google.com",
			userID:      testSvcUserID,
			setupMocks: func(ctx context.Context, repo *repoMocks.Repository, keygen *keygenMocks.KeyGenerator) {
				// 1. Build expected input model passed to Repo
				// Note: Code is empty because it is commented out in your service
				expectedInput := &model.Bookmark{
					Description: "Google Search",
					URL:         "https://google.com",
					UserID:      fixtures.GetUUID(testSvcUserID),
				}

				// 2. Prepare mocked repository output record
				mockedSavedBookmark := &model.Bookmark{
					Base: model.Base{
						ID: uuid.New(),
					},
					Description: "Google Search",
					URL:         "https://google.com",
					UserID:      fixtures.GetUUID(testSvcUserID),
				}

				// 3. Mock the repository call
				repo.On("CreateBookmark", ctx, expectedInput).Return(mockedSavedBookmark, nil)
			},
			expectedOut: &model.Bookmark{
				Description: "Google Search",
				URL:         "https://google.com",
			},
			expectedErr: nil,
		},
		{
			name:        "Error - Repository Failure",
			description: "Google Search",
			url:         "https://google.com",
			userID:      testSvcUserID,
			setupMocks: func(ctx context.Context, repo *repoMocks.Repository, keygen *keygenMocks.KeyGenerator) {
				expectedInput := &model.Bookmark{
					Description: "Google Search",
					URL:         "https://google.com",
					UserID:      fixtures.GetUUID(testSvcUserID),
				}

				// Mock the repository call to return an error
				repo.On("CreateBookmark", ctx, expectedInput).
					Return((*model.Bookmark)(nil), errors.New("database connectivity error"))
			},
			expectedOut: nil,
			expectedErr: errors.New("database connectivity error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			// Initialize mock controllers

			mockRepo := repoMocks.NewRepository(t)
			mockKeygen := keygenMocks.NewKeyGenerator(t) // Adjust constructor name to match your generated mock setup

			tc.setupMocks(ctx, mockRepo, mockKeygen)

			// Construct service injecting the mocked dependencies
			service := &bookmarkService{
				bookmarkRepo: mockRepo,
				keygen:       mockKeygen,
			}

			output, err := service.AddBookmark(ctx, tc.description, tc.url, tc.userID)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Nil(t, output)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.Equal(t, tc.expectedOut.Description, output.Description)
			assert.Equal(t, tc.expectedOut.URL, output.URL)
			assert.Equal(t, tc.expectedOut.Code, output.Code)
			assert.Equal(t, fixtures.GetUUID(tc.userID), output.UserID)
		})
	}
}

package bookmark

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	repoMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/bookmark/mocks"
)

const (
	testSvcGetUserID = "133b3b42-70b9-456c-82e7-bf1b570e6c51"
)

func TestBookmarkService_GetBookmarks(t *testing.T) {
	t.Parallel()

	mockedBookmarksList := []*model.Bookmark{
		{
			Description: "Google Search",
			URL:         "https://google.com",
			Code:        "12345abc",
		},
	}
	testCases := []struct {
		name           string
		userID         string
		page           int
		limit          int
		setupRepo      func(ctx context.Context, repo *repoMocks.Repository)
		expectedResult *BookmarkResult
		expectedError  error
	}{
		{
			name:   "Success - Get Bookmarks Page 1",
			userID: testSvcGetUserID,
			page:   1, // (1-1) * 10 = 0 offset
			limit:  10,
			setupRepo: func(ctx context.Context, repo *repoMocks.Repository) {
				// 1. Verify GetBookmarks receives calculation offset '0'
				repo.On("GetBookmarks", ctx, testSvcGetUserID, 0, 10).
					Return(mockedBookmarksList, nil)

				// 2. Verify subsequent count retrieval call
				repo.On("GetBookmarkCount", ctx, testSvcGetUserID).
					Return(int64(1), nil)
			},
			expectedResult: &BookmarkResult{
				Bookmarks: mockedBookmarksList,
				Count:     1,
			},
			expectedError: nil,
		},
		{
			name:   "Success - Get Bookmarks Page 3 Offset Calculation Check",
			userID: testSvcGetUserID,
			page:   3, // (3-1) * 5 = 10 offset
			limit:  5,
			setupRepo: func(ctx context.Context, repo *repoMocks.Repository) {
				repo.On("GetBookmarks", ctx, testSvcGetUserID, 10, 5).
					Return(mockedBookmarksList, nil)

				repo.On("GetBookmarkCount", ctx, testSvcGetUserID).
					Return(int64(11), nil)
			},
			expectedResult: &BookmarkResult{
				Bookmarks: mockedBookmarksList,
				Count:     11,
			},
			expectedError: nil,
		},
		{
			name:   "Error - Main List Fetch Failure",
			userID: testSvcGetUserID,
			page:   1,
			limit:  10,
			setupRepo: func(ctx context.Context, repo *repoMocks.Repository) {
				repo.On("GetBookmarks", ctx, testSvcGetUserID, 0, 10).
					Return(([]*model.Bookmark)(nil), errors.New("list pull db crash"))
			},
			expectedResult: nil,
			expectedError:  errors.New("list pull db crash"),
		},
		{
			name:   "Error - Total Counter Fetch Failure",
			userID: testSvcGetUserID,
			page:   1,
			limit:  10,
			setupRepo: func(ctx context.Context, repo *repoMocks.Repository) {
				repo.On("GetBookmarks", ctx, testSvcGetUserID, 0, 10).
					Return(mockedBookmarksList, nil)

				repo.On("GetBookmarkCount", ctx, testSvcGetUserID).
					Return(int64(0), errors.New("count aggregate db crash"))
			},
			expectedResult: nil,
			expectedError:  errors.New("count aggregate db crash"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			// Initialize mock repository layer controller
			mockRepo := repoMocks.NewRepository(t)
			tc.setupRepo(ctx, mockRepo)

			// Construct service injecting the mocked repository configuration
			service := &bookmarkService{
				bookmarkRepo: mockRepo,
			}

			output, err := service.GetBookmarks(ctx, tc.userID, tc.page, tc.limit)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, output)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.Equal(t, tc.expectedResult.Count, output.Count)
			assert.Equal(t, len(tc.expectedResult.Bookmarks), len(output.Bookmarks))
		})
	}
}

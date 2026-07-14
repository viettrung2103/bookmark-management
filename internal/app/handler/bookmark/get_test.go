package bookmark

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	bookmarkService "github.com/viettrung2103/bookmark-management/internal/app/service/bookmark"
	bookmarkMock "github.com/viettrung2103/bookmark-management/internal/app/service/bookmark/mocks"
)

const (
	testGetUserID = "12345678-1234-1234-1234-123456789012"
)

func TestBookmarkHandler_GetBookmarks(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *bookmarkMock.Service

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success - get list of bookmarks with pagination",
			setupRequest: func(ctx *gin.Context) {
				ctx.Set("claims", newTestClaims(testGetUserID, "testuser"))
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=10", nil)
			},
			setupMockService: func(ctx context.Context) *bookmarkMock.Service {
				serviceMock := bookmarkMock.NewService(t)

				mockResult := &bookmarkService.BookmarkResult{
					Bookmarks: []*model.Bookmark{
						{
							Description: "Google",
							URL:         "https://google.com",
							Code:        "abc1",
						},
					},
					Count: 1,
				}

				serviceMock.On("GetBookmarks", mock.Anything, testGetUserID, 1, 10).
					Return(mockResult, nil)

				return serviceMock
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"data":[{"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","description":"Google","id":"00000000-0000-0000-0000-000000000000","url":"https://google.com","code":"abc1"}],"pagination":{"count":1,"page":1,"limit":10}}`,
		},
		{
			name: "error - service internal failure",
			setupRequest: func(ctx *gin.Context) {
				ctx.Set("claims", newTestClaims(testGetUserID, "testuser"))
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=20", nil)
			},
			setupMockService: func(ctx context.Context) *bookmarkMock.Service {
				serviceMock := bookmarkMock.NewService(t)

				// 🟢 Fix: Return a nil pointer cast to the correct type along with the error
				serviceMock.On("GetBookmarks", mock.Anything, testGetUserID, 1, 20).
					Return((*bookmarkService.BookmarkResult)(nil), errors.New("database read crash"))

				return serviceMock
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"message":"Internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockService(ctx)

			testHandler := NewHandler(mockSvc)
			testHandler.GetBookmarks(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedResponse != "" {
				assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
			}
		})
	}
}

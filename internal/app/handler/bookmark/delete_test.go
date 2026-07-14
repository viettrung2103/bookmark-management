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
	bookmarkMock "github.com/viettrung2103/bookmark-management/internal/app/service/bookmark/mocks"
)

const (
	// Reusing the same strategy: clean global constants for cross-verifying arguments
	testDeleteUserID     = "12345678-1234-1234-1234-123456789012"
	testDeleteBookmarkID = "87654321-4321-4321-4321-210987654321"
)

func TestBookmarkHandler_DeleteBookmark(t *testing.T) {
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
			name: "success - delete bookmark",
			setupRequest: func(ctx *gin.Context) {
				// 1. Inject the claims safely to satisfy requestutils.GetUserIDFromRequest
				ctx.Set("claims", newTestClaims(testDeleteUserID, "testuser"))

				// 2. Setup the target mock request and path parameters
				ctx.Request = httptest.NewRequest(http.MethodDelete, "/v1/bookmarks/"+testDeleteBookmarkID, nil)
				ctx.Params = gin.Params{
					{Key: "id", Value: testDeleteBookmarkID},
				}
			},
			setupMockService: func(ctx context.Context) *bookmarkMock.Service {
				serviceMock := bookmarkMock.NewService(t)

				// Service layer expects the Context, the verified User ID, and the Target Bookmark ID
				serviceMock.On("DeleteBookmarkByID", mock.Anything, testDeleteUserID, testDeleteBookmarkID).
					Return(nil)

				return serviceMock
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"Success"}`,
		},
		{
			name: "error - service internal failure",
			setupRequest: func(ctx *gin.Context) {
				ctx.Set("claims", newTestClaims(testDeleteUserID, "testuser"))

				ctx.Request = httptest.NewRequest(http.MethodDelete, "/v1/bookmarks/"+testDeleteBookmarkID, nil)
				ctx.Params = gin.Params{
					{Key: "id", Value: testDeleteBookmarkID},
				}
			},
			setupMockService: func(ctx context.Context) *bookmarkMock.Service {
				serviceMock := bookmarkMock.NewService(t)

				serviceMock.On("DeleteBookmarkByID", mock.Anything, testDeleteUserID, testDeleteBookmarkID).
					Return(errors.New("failed to delete from database"))

				return serviceMock
			},
			expectedStatus: http.StatusInternalServerError,
			// Matches your explicit global application configuration response wrapper: response.InternalErrResponse
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
			testHandler.DeleteBookmark(ctx)

			// Assertions
			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedResponse != "" {
				assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
			}
		})
	}
}

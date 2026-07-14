package bookmark

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	bookmarkMock "github.com/viettrung2103/bookmark-management/internal/app/service/bookmark/mocks"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
)

const (
	testEditUserID     = "12345678-1234-1234-1234-123456789012"
	testEditBookmarkID = "87654321-4321-4321-4321-210987654321"
)

func TestBookmarkHandler_EditBookmark(t *testing.T) {
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
			name: "success - edit bookmark",
			setupRequest: func(ctx *gin.Context) {
				// 1. Inject auth claims
				ctx.Set("claims", newTestClaims(testEditUserID, "testuser"))

				// 2. Setup JSON payload
				body := map[string]any{
					"description": "Updated Google Search",
					"url":         "https://google.com",
				}
				jsonBody, _ := json.Marshal(body)

				// 3. Mount parameters and request
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+testEditBookmarkID, bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Params = gin.Params{
					{Key: "id", Value: testEditBookmarkID},
				}
			},
			setupMockService: func(ctx context.Context) *bookmarkMock.Service {
				serviceMock := bookmarkMock.NewService(t)

				// Expect service tracking arguments matching parameters passed from controller
				serviceMock.On("EditBookmarkByID", ctx, testEditUserID, testEditBookmarkID, "Updated Google Search", "https://google.com").
					Return(nil)

				return serviceMock
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"Success"}`,
		},
		{
			name: "error - bookmark duplication conflict",
			setupRequest: func(ctx *gin.Context) {
				ctx.Set("claims", newTestClaims(testEditUserID, "testuser"))

				body := map[string]any{
					"description": "Duplicate URL",
					"url":         "https://duplicate.com",
				}
				jsonBody, _ := json.Marshal(body)

				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+testEditBookmarkID, bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Params = gin.Params{
					{Key: "id", Value: testEditBookmarkID},
				}
			},
			setupMockService: func(ctx context.Context) *bookmarkMock.Service {
				serviceMock := bookmarkMock.NewService(t)

				// Injecting dbutils.ErrDuplication to target your specific switch case matching rule
				serviceMock.On("EditBookmarkByID", ctx, testEditUserID, testEditBookmarkID, "Duplicate URL", "https://duplicate.com").
					Return(dbutils.ErrDuplication)

				return serviceMock
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"duplicated column"}`, // Adjust error message match to equal dbutils.ErrDuplication.Error()
		},
		{
			name: "error - service internal failure",
			setupRequest: func(ctx *gin.Context) {
				ctx.Set("claims", newTestClaims(testEditUserID, "testuser"))

				body := map[string]any{
					"description": "Broken Request",
					"url":         "https://error.com",
				}
				jsonBody, _ := json.Marshal(body)

				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+testEditBookmarkID, bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
				ctx.Params = gin.Params{
					{Key: "id", Value: testEditBookmarkID},
				}
			},
			setupMockService: func(ctx context.Context) *bookmarkMock.Service {
				serviceMock := bookmarkMock.NewService(t)

				serviceMock.On("EditBookmarkByID", ctx, testEditUserID, testEditBookmarkID, "Broken Request", "https://error.com").
					Return(errors.New("unexpected database disconnection"))

				return serviceMock
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"message":"Internal server error"}`, // Matches your global response.InternalErrResponse fields
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
			testHandler.EditBookmark(ctx)

			// Assertions
			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedResponse != "" {
				assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
			}
		})
	}
}

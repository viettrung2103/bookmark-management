package bookmark

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	//"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	bookmarkMock "github.com/viettrung2103/bookmark-management/internal/app/service/bookmark/mocks"
)

const (
	// Shared Test Constant to align the Handler and Service Mock expectations cleanly
	testUserID = "12345678-1234-1234-1234-123456789012"
)

// newTestClaims generates a reusable jwt.MapClaims layout mimicking ToMapClaim production logic
func newTestClaims(uid string, username string) jwt.MapClaims {
	return jwt.MapClaims{
		"uid":      uid,
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
}

func TestBookmarkHandler_AddBookmark(t *testing.T) {
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
			name: "success - create bookmark",
			setupRequest: func(ctx *gin.Context) {
				// Inject the claims effortlessly via the shared factory function
				ctx.Set("claims", newTestClaims(testUserID, "testuser"))

				body := map[string]any{
					"description": "Google Search",
					"url":         "https://www.google.com",
				}
				jsonBody, _ := json.Marshal(body)

				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockService: func(ctx context.Context) *bookmarkMock.Service {
				serviceMock := bookmarkMock.NewService(t)

				testUUID := "87654321-4321-4321-4321-210987654321"
				mockBookmark := &model.Bookmark{
					Base: model.Base{
						ID: uuid.MustParse(testUUID),
					},
					Description: "Google Search",
					URL:         "https://www.google.com",
					Code:        "abc1235",
				}

				serviceMock.On("AddBookmark", mock.Anything, "Google Search", "https://www.google.com", testUserID).
					Return(mockBookmark, nil)

				return serviceMock
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"data":{"id":"87654321-4321-4321-4321-210987654321","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","description":"Google Search","url":"https://www.google.com","code":"abc1235"},"message":"Create a bookmark successfully"}`,
		},
		{
			name: "error - service internal failure",
			setupRequest: func(ctx *gin.Context) {
				// Reusing the same factory helper for the error test track
				ctx.Set("claims", newTestClaims(testUserID, "testuser"))

				body := map[string]any{
					"description": "Google Search",
					"url":         "https://www.google.com",
				}
				jsonBody, _ := json.Marshal(body)

				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockService: func(ctx context.Context) *bookmarkMock.Service {
				serviceMock := bookmarkMock.NewService(t)

				serviceMock.On("AddBookmark", mock.Anything, "Google Search", "https://www.google.com", testUserID).
					Return((*model.Bookmark)(nil), errors.New("database connection down"))

				return serviceMock
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: ``,
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
			testHandler.AddBookmark(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedResponse != "" {
				assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
			}
		})
	}
}

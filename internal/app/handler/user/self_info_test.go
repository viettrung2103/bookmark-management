package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	userMock "github.com/viettrung2103/bookmark-management/internal/app/service/user/mocks"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	//"gorm.io/gorm"
)

// TestUserHandler_SelfInfo tests the SelfInfo method of the UserHandler
func TestUserHandler_SelfInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *userMock.UserService

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)

				// Change "userId" to whatever key your middleware actually uses (e.g., "user_id", "userID", etc.)
				mockClaims := jwt.MapClaims{"uid": "user-123"}
				ctx.Set("claims", mockClaims)
			},
			setupMockService: func(ctx context.Context) *userMock.UserService {
				serviceMock := userMock.NewUserService(t)
				testUUID := "12345678-1234-1234-1234-123456789012"

				mockUser := &model.User{
					Base: model.Base{
						ID: uuid.MustParse(testUUID),
					},
					Username:    "testuser",
					Email:       "test@example.com",
					DisplayName: "testuser",
				}

				serviceMock.On("SelfInfo", ctx, "user-123").Return(mockUser, nil)
				return serviceMock
			},
			expectedStatus: http.StatusOK,
			// Notice this expects the raw user object because your handler does `c.JSON(http.StatusOK, user)`
			expectedResponse: `{"data":{"id":"12345678-1234-1234-1234-123456789012","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","display_name":"testuser","username":"testuser","email":"test@example.com"}}`,
		},
		{
			name: "user not found",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				mockClaims := jwt.MapClaims{"uid": "missing-user-999"}
				ctx.Set("claims", mockClaims)
			},
			setupMockService: func(ctx context.Context) *userMock.UserService {
				serviceMock := userMock.NewUserService(t)
				// Return GORM's record not found error
				serviceMock.On("SelfInfo", ctx, "missing-user-999").Return((*model.User)(nil), dbutils.ErrRecordNotFound)
				return serviceMock
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"User Not Found"}`,
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)

				// 1. Create the MapClaims with the standard user ID
				mockClaims := jwt.MapClaims{
					"uid": "user-123",
				}
				// 2. Set it in the context under the "claims" key
				ctx.Set("claims", mockClaims)
			},
			setupMockService: func(ctx context.Context) *userMock.UserService {
				serviceMock := userMock.NewUserService(t)
				serviceMock.On("SelfInfo", ctx, "user-123").Return((*model.User)(nil), errors.New("database connection lost"))
				return serviceMock
			},
			expectedStatus: http.StatusInternalServerError,
			// IMPORTANT: You need to replace this expected response with exactly what your `response.InternalErrResponse` serializes to.
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
			testHandler.SelfInfo(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Safely compare the JSON output
			if tc.expectedResponse != "" {
				assert.Equal(t, tc.expectedResponse, rec.Body.String())
			}
		})
	}
}

// TestUserHandler_EditSelfInfo tests the EditSelfInfo method of the UserHandler
func TestUserHandler_EditSelfInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *userMock.UserService

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",
			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"display_name": "Updated Name",
					"email":        "updated@example.com",
				}
				jsonBody, _ := json.Marshal(body)

				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")

				// SIMULATING THE JWT MIDDLEWARE: Inject the user ID into the context
				mockClaims := jwt.MapClaims{"uid": "user-123"}
				ctx.Set("claims", mockClaims)
				//ctx.Set("userId", "user-123") // Make sure "userId" matches what GetUserIDFromRequest expects
			},
			setupMockService: func(ctx context.Context) *userMock.UserService {
				serviceMock := userMock.NewUserService(t)
				// Expect the exact parameters from the request and return nil error
				serviceMock.On("EditInfoByID", ctx, "user-123", "Updated Name", "updated@example.com").Return(nil)
				return serviceMock
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"Edit current user successfully!"}`,
		},
		{
			name: "duplication error (email already exists)",
			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"display_name": "Updated Name",
					"email":        "duplicate@example.com",
				}
				jsonBody, _ := json.Marshal(body)

				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
				mockClaims := jwt.MapClaims{"uid": "user-123"}
				ctx.Set("claims", mockClaims)
			},
			setupMockService: func(ctx context.Context) *userMock.UserService {
				serviceMock := userMock.NewUserService(t)
				// Simulate the DB rejecting the update due to a duplicate email
				serviceMock.On("EditInfoByID", ctx, "user-123", "Updated Name", "duplicate@example.com").Return(dbutils.ErrDuplication)
				return serviceMock
			},
			expectedStatus: http.StatusBadRequest,
			// Assuming dbutils.ErrDuplication.Error() returns something like "record duplicated"
			expectedResponse: `{"message":"` + dbutils.ErrDuplication.Error() + `"}`,
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"display_name": "Updated Name",
					"email":        "updated@example.com",
				}
				jsonBody, _ := json.Marshal(body)

				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
				mockClaims := jwt.MapClaims{"uid": "user-123"}
				ctx.Set("claims", mockClaims)
			},
			setupMockService: func(ctx context.Context) *userMock.UserService {
				serviceMock := userMock.NewUserService(t)
				serviceMock.On("EditInfoByID", ctx, "user-123", "Updated Name", "updated@example.com").Return(errors.New("db down"))
				return serviceMock
			},
			expectedStatus: http.StatusInternalServerError,
			// You may need to adjust this string to match what response.InternalErrResponse outputs
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
			testHandler.EditSelfInfo(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedResponse != "" {
				assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
			}
		})
	}
}

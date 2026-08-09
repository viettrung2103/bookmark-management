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
	"github.com/stretchr/testify/assert"
	userMock "github.com/viettrung2103/bookmark-management/internal/app/service/user/mocks"

	//userMocks"github.com/viettrung2103/bookmark-management/internal/app/service/user/mocks"
	"github.com/viettrung2103/bookmark-management/internal/app/service/user"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
)

// TestUserHandler_Login tests the Login function
func TestUserHandler_Login(t *testing.T) {
	t.Parallel()
	//cfg, err := config.NewConfig()
	//if err != nil {
	//	panic(err)
	//}

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *userMock.Service

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",
			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"username": "testuser",
					"password": "validpassword123",
				}
				jsonBody, _ := json.Marshal(body)

				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewReader(jsonBody))
				// CRITICAL: Gin requires the Content-Type header to bind JSON properly
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockService: func(ctx context.Context) *userMock.Service {
				serviceMock := userMock.NewService(t)
				// Return a fake token and no error
				serviceMock.On("Login", ctx, "testuser", "validpassword123").Return("fake-jwt-token-123", nil)
				return serviceMock
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"data":"fake-jwt-token-123","message":"Logged in successfully!"}`,
		},
		{
			name: "invalid credentials error",
			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"username": "testuser",
					"password": "wrongpassword",
				}
				jsonBody, _ := json.Marshal(body)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockService: func(ctx context.Context) *userMock.Service {
				serviceMock := userMock.NewService(t)
				// Simulate the service returning ErrInvalidCreditials
				serviceMock.On("Login", ctx, "testuser", "wrongpassword").Return("", user.ErrInvalidCreditials)
				return serviceMock
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid username or password"}`,
		},
		{
			name: "record not found error",
			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"username": "nonexistentuser",
					"password": "somepassword",
				}
				jsonBody, _ := json.Marshal(body)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockService: func(ctx context.Context) *userMock.Service {
				serviceMock := userMock.NewService(t)
				// Simulate the DB not finding the user
				serviceMock.On("Login", ctx, "nonexistentuser", "somepassword").Return("", dbutils.ErrRecordNotFound)
				return serviceMock
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid username or password"}`,
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"username": "testuser",
					"password": "validpassword123",
				}
				jsonBody, _ := json.Marshal(body)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockService: func(ctx context.Context) *userMock.Service {
				serviceMock := userMock.NewService(t)
				// Simulate an unexpected error (like a DB connection failure)
				serviceMock.On("Login", ctx, "testuser", "validpassword123").Return("", errors.New("database connection refused"))
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
			testHandler.Login(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

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
	"github.com/stretchr/testify/mock"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/app/service/mocks"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
)

func TestUserHandler_Register(t *testing.T) {
	t.Parallel()
	//cfg, err := config.NewConfig()
	//if err != nil {
	//	panic(err)
	//}

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.UserService

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",
			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"display_name": "Test User",
					"username":     "testuser",
					"password":     "validpassword123", // Must be > 8
					"email":        "test@example.com",
				}
				jsonBody, _ := json.Marshal(body)

				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockService: func(ctx context.Context) *mocks.UserService {
				serviceMock := mocks.NewUserService(t)

				// Create the expected return model
				mockUser := &model.User{
					ID:       "user-123",
					Username: "testuser",
					Email:    "test@example.com",
				}

				serviceMock.On("CreateUser", mock.Anything, "Test User", "testuser", "validpassword123", "test@example.com").Return(mockUser, nil)
				return serviceMock
			},
			expectedStatus: http.StatusOK,
			// Notice how data aligns with the model returned by the mock
			expectedResponse: `{"data":{"id":"user-123","display_name":"","username":"testuser","email":"test@example.com"},"message":"Register an user successfully"}`,
		},
		{
			name: "duplication error",
			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"display_name": "Test User",
					"username":     "duplicateuser",
					"password":     "validpassword123",
					"email":        "dup@example.com",
				}
				jsonBody, _ := json.Marshal(body)

				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockService: func(ctx context.Context) *mocks.UserService {
				serviceMock := mocks.NewUserService(t)
				// Return the dbutils.ErrDuplication error
				serviceMock.On("CreateUser", mock.Anything, "Test User", "duplicateuser", "validpassword123", "dup@example.com").Return((*model.User)(nil), dbutils.ErrDuplication)
				return serviceMock
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"User or Email already exist"}`,
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"display_name": "Test User",
					"username":     "testuser",
					"password":     "validpassword123",
					"email":        "test@example.com",
				}
				jsonBody, _ := json.Marshal(body)

				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockService: func(ctx context.Context) *mocks.UserService {
				serviceMock := mocks.NewUserService(t)
				// Return a generic error
				serviceMock.On("CreateUser", mock.Anything, "Test User", "testuser", "validpassword123", "test@example.com").Return((*model.User)(nil), errors.New("database connection lost"))
				return serviceMock
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"message":"Internal server error"}`,
		},
		//
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockService(ctx)
			testHandler := NewHandler(mockSvc)
			testHandler.Register(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

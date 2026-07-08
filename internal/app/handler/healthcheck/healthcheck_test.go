package healthcheck

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/viettrung2103/bookmark-management/internal/app/service/mocks"
)

// TestShortenLinkHandler tests the ShortenLinkHandler function
func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.HealthCheckService

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)

			},
			setupMockService: func(ctx context.Context) *mocks.HealthCheckService {
				serviceMock := mocks.NewHealthCheckService(t)
				serviceMock.On("HealthCheck", mock.Anything).Return(nil)
				return serviceMock
			},

			expectedStatus:   http.StatusOK,
			expectedResponse: `{"redis":"reachable","status":"UP"}`,
		},
		{
			name: "redis is down",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockService: func(ctx context.Context) *mocks.HealthCheckService {
				serviceMock := mocks.NewHealthCheckService(t)

				// 1. Return an actual error here to trigger the `if err != nil` block in your handler
				// Note: If your HealthCheck interface returns a struct AND an error,
				// you would write something like: Return(nil, errors.New("redis error"))
				serviceMock.On("HealthCheck", mock.Anything).Return(errors.New("redis connection refused"))

				return serviceMock
			},
			// 2. Expect the 503 Service Unavailable status
			expectedStatus: http.StatusServiceUnavailable,

			// 3. Expect the JSON payload defined in your handler
			expectedResponse: `{"error":"Service Unavailable","redis":"unreachable","status":"DOWN"}`,
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
			testHandler.HealthCheck(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

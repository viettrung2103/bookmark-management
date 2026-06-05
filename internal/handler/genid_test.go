package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/internal/service/mocks"
)

var testErr = errors.New("test error")
var testUUID = "123e4567-e89b-12d3-a456-426614174000"

func TestGenPassHandler_GeneratePassword(t *testing.T) {
	t.Parallel()
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.GenId

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "v1/links//health-check", nil)
			},
			setupMockService: func(ctx context.Context) *mocks.GenId {
				serviceMock := mocks.NewGenId(t)
				serviceMock.On("GenerateId").Return(testUUID, nil)
				return serviceMock
			},

			expectedStatus:   http.StatusOK,
			expectedResponse: `{"instance_id":"123e4567-e89b-12d3-a456-426614174000","message":"OK","service_name":"bookmark_service"}`,
		},
		{
			name: "service failed",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "v1/links//health-check", nil)
			},
			setupMockService: func(ctx context.Context) *mocks.GenId {
				serviceMock := mocks.NewGenId(t)
				serviceMock.On("GenerateId").Return("", testErr)
				return serviceMock
			},

			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"error":"Internal Server Err"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockService(ctx)
			testHandler := NewGenId(mockSvc, cfg)
			testHandler.GenerateId(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

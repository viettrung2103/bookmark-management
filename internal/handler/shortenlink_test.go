package handler

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
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/internal/service/mocks"
)

// var testErr = errors.New("test error")
var testCode = "abc1235"
var basePath = "/v1/links"

func TestShortenLinkHandler(t *testing.T) {
	t.Parallel()
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.ShortenUrl

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"exp": 10,
					"url": "https://www.google.com",
				}

				// 2. Convert the map to JSON bytes
				jsonBody, err := json.Marshal(body)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}

				ctx.Request = httptest.NewRequest(http.MethodPost, "v1/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("ShortenUrlWithExpiringTime", mock.Anything, "https://www.google.com", 10).Return(testCode, nil)
				return serviceMock
			},

			expectedStatus:   http.StatusOK,
			expectedResponse: `{"code":"abc1235","message":"Shorten URL generated successfully"}`,
		},
		{
			name: "wrong request body",

			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"exp": 10,
					//"url": "https://www.google.com",
				}

				// 2. Convert the map to JSON bytes
				jsonBody, err := json.Marshal(body)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}

				ctx.Request = httptest.NewRequest(http.MethodPost, "v1/links//shorten", bytes.NewReader(jsonBody))
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("ShortenUrlWithExpiringTime", mock.Anything, "", 10).Return("", errors.New("Internal Server Err"))
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
			testHandler := NewShortenLink(mockSvc, cfg)
			testHandler.ShortenUrlLink(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

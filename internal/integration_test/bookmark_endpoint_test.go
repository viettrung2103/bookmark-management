package intergration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/api"
	"github.com/viettrung2103/bookmark-management/internal/config"
	pkgRedis "github.com/viettrung2103/bookmark-management/pkg/redis"
)

func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "normal case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("GET", "/health-check", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},

			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"redis":`,
		},
		{
			name: "wrong endpoint",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("POST", "/health-check", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},

			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			redisMocks := pkgRedis.InitMockRedis(t)

			testApi := api.NewEngine(&config.Config{}, redisMocks)
			recorder := tc.setupTestHTTP(testApi)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}

func TestShortenUrlEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{

			name: "normal case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"exp": 10,
					"url": "https://www.google.com",
				}
				// 2. Convert the map to JSON bytes
				jsonBody, err := json.Marshal(body)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}
				req, _ := http.NewRequest("POST", "/v1/links/shorten", bytes.NewReader(jsonBody))
				req.Header.Set("Content-Type", "application/json")

				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},

			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"code":`,
		},
		{
			name: "wrong endpoint",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("GET", "/v1/links/shorten", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},

			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			redisMocks := pkgRedis.InitMockRedis(t)

			testApi := api.NewEngine(&config.Config{}, redisMocks)
			recorder := tc.setupTestHTTP(testApi)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}

const testExpTime = 1000

func TestRedirectEndpoint(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	testCases := []struct {
		name string

		// simulate having code in redis
		setupMockRedis func() *redis.Client
		setupTestHTTP  func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode int
		expectedUrl        string
	}{
		{
			name: "normal case",
			setupMockRedis: func() *redis.Client {
				mock := pkgRedis.InitMockRedis(t)
				mock.Set(ctx, "1234567", "https://test.com", testExpTime)
				//mock.On("StoreUrl", "1234567", "https://test.com", testExpTime)
				return mock
			},

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("GET", "/v1/links/redirect/1234567", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},

			expectedStatusCode: http.StatusFound,
			expectedUrl:        `https://test.com`,
		},
		{
			name: "wrong endpoint",
			setupMockRedis: func() *redis.Client {
				mock := pkgRedis.InitMockRedis(t)
				mock.Set(ctx, "1234567", "https://test.com", testExpTime)
				//mock.On("StoreUrl", "1234567", "https://test.com", testExpTime)
				return mock
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("POST", "/v1/links/redirect/1234567", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},

			expectedStatusCode: http.StatusNotFound,
			expectedUrl:        ``,
		},
		{
			name: "wrong code",
			setupMockRedis: func() *redis.Client {
				mock := pkgRedis.InitMockRedis(t)
				mock.Set(ctx, "1234567", "https://test.com", testExpTime)
				//mock.On("StoreUrl", "1234567", "https://test.com", testExpTime)
				return mock
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest("POST", "/v1/links/redirect/234567", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},

			expectedStatusCode: http.StatusNotFound,
			expectedUrl:        ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			redisMocks := tc.setupMockRedis()

			testApi := api.NewEngine(&config.Config{}, redisMocks)
			recorder := tc.setupTestHTTP(testApi)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Equal(t, tc.expectedUrl, recorder.Header().Get("Location"))

		})
	}
}

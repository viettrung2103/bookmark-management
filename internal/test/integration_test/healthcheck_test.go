package intergration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/api"
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	pkgRedis "github.com/viettrung2103/bookmark-management/pkg/redis"
	"gorm.io/gorm"
)

// TestHealthCheckEndpoint tests the health check endpoint
func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder
		setupDB       func() *gorm.DB

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
			setupDB: func() *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
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
			setupDB: func() *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},

			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			redisMocks := pkgRedis.InitMockRedis(t)
			//engine := gin.New()
			// generate test cache
			fixtures := tc.setupDB()

			testApi := api.NewEngine(&api.EngineOpts{
				Engine: gin.New(),
				Cfg:    &config.Config{},
				Redis:  redisMocks,
				SqlDB:  fixtures,
			})

			//engine, &config.Config{}, redisMocks, fixtures)
			recorder := tc.setupTestHTTP(testApi)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}

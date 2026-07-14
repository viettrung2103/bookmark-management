package url

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	link "github.com/viettrung2103/bookmark-management/internal/app/service/url"
	urlMock "github.com/viettrung2103/bookmark-management/internal/app/service/url/mocks"
	"github.com/viettrung2103/bookmark-management/internal/config"
)

func TestShortenLinkHandler_RedirectUrl(t *testing.T) {
	t.Parallel()
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		name string

		setupRequest func(c *gin.Context)
		setupMockSvc func(ctx context.Context) *urlMock.URLService

		expectedStatus   int
		expectedUrl      string
		expectedResponse string // Added to validate JSON bodies on errors
	}{
		{
			name: "success - redirect to original url",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/1234567", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "1234567"}}
			},
			setupMockSvc: func(ctx context.Context) *urlMock.URLService {
				serviceMock := urlMock.NewURLService(t)
				serviceMock.On("GetLinkFromCode", ctx, "1234567").Return("https://google.com", nil)
				return serviceMock
			},
			expectedStatus: http.StatusFound,
			expectedUrl:    "https://google.com",
		},
		{
			name: "error - missing code parameter",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/", nil)
				ctx.Params = gin.Params{{Key: "code", Value: ""}}
			},
			setupMockSvc: func(ctx context.Context) *urlMock.URLService {
				// Service should not even be called when input parameters are missing
				return urlMock.NewURLService(t)
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"code is required"}`,
		},
		{
			name: "error - code does not exist",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/invalidcode", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "invalidcode"}}
			},
			setupMockSvc: func(ctx context.Context) *urlMock.URLService {
				serviceMock := urlMock.NewURLService(t)
				// Simulating link.ErrCodeDoesNotExist from service layer
				serviceMock.On("GetLinkFromCode", ctx, "invalidcode").Return("", link.ErrCodeDoesNotExist)
				return serviceMock
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"url not found"}`,
		},
		{
			name: "error - internal server error from database",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/dberrorcode", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "dberrorcode"}}
			},
			setupMockSvc: func(ctx context.Context) *urlMock.URLService {
				serviceMock := urlMock.NewURLService(t)
				// Simulating a standard connection or querying failure
				serviceMock.On("GetLinkFromCode", ctx, "dberrorcode").Return("", errors.New("database connection down"))
				return serviceMock
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"error":"Internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockSvc(ctx)
			testHandler := NewShortenLink(mockSvc, cfg)
			testHandler.RedirectUrl(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedUrl != "" {
				assert.Equal(t, tc.expectedUrl, rec.Header().Get("Location"))
			}

			if tc.expectedResponse != "" {
				assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
			}
		})
	}
}

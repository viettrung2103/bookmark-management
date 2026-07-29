package intergration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/api"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	"github.com/viettrung2103/bookmark-management/pkg/jwtutils"
	jwtMocks "github.com/viettrung2103/bookmark-management/pkg/jwtutils/mocks"
	pkgRedis "github.com/viettrung2103/bookmark-management/pkg/redis"
	"gorm.io/gorm"
)

const (
	mockAuthUserID     = "133b3b42-70b9-456c-82e7-bf1b570e6c51"
	testUsername       = "johndoe99"
	mockAuthUserClaims = "mocked-auth-user-claims"
	testTokenHeader    = "Bearer test-valid-jwt-token"
)

// A generic structure to capture standard API response formats if needed

func TestEndpoint_Bookmark_Create(t *testing.T) {
	t.Parallel()

	// Seed tracking identifiers
	//targetBookmarkUUID := uuid.New()

	testCases := []struct {
		name                 string
		setupTestHTTP        func(api api.Engine) *httptest.ResponseRecorder
		setupDB              func() *gorm.DB
		setupMockJwt         func(mockVal *jwtMocks.JWTValidator)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success - Create Bookmark",
			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(
					http.MethodPost,
					"/v1/bookmarks",
					bytes.NewBuffer([]byte(`{
						"description": "Go Packages",
						"url": "https://pkg.go.dev"
					}`)))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", testTokenHeader)

				recorder := httptest.NewRecorder()
				app.ServeHTTP(recorder, req)
				return recorder
			},
			setupDB: func() *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupMockJwt: func(mockVal *jwtMocks.JWTValidator) {
				// Import your JWT package if needed (e.g., "github.com/golang-jwt/jwt/v5")
				// active claims map matching what your middleware reads
				mockClaims := jwtutils.GetMapClaim(mockAuthUserID, testUsername)

				mockVal.On("ValidateJWT", "test-valid-jwt-token").Return(mockClaims, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `"url":"https://pkg.go.dev"`,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			redisMocks := pkgRedis.InitMockRedis(t)

			// 1. Reset database fixtures
			db := tc.setupDB()

			// 2. Build mock validator for privateBase protection layers
			mockJwtValidator := jwtMocks.NewJWTValidator(t)
			tc.setupMockJwt(mockJwtValidator)

			// 3. Mount engine with mock dependencies injected
			testRouter := api.NewEngine(&api.EngineOpts{
				Engine: gin.New(),
				Cfg:    &config.Config{},
				Redis:  redisMocks,
				SqlDB:  db,
				JwtGen: nil, // Setup not required for standard resource routes
				JwtVal: mockJwtValidator,
			})

			// 4. Fire the test execution
			record := tc.setupTestHTTP(testRouter)

			// 5. Run validation assertions
			assert.Equal(t, tc.expectedStatusCode, record.Code)
			if tc.expectedResponseBody != "" {
				assert.Contains(t, record.Body.String(), tc.expectedResponseBody)
			}
		})
	}
}

func TestEndpoint_Bookmark_Get(t *testing.T) {
	t.Parallel()

	// Seed tracking identifiers
	targetBookmarkUUID := uuid.New()

	testCases := []struct {
		name                 string
		setupTestHTTP        func(api api.Engine) *httptest.ResponseRecorder
		setupDB              func() *gorm.DB
		setupMockJwt         func(mockVal *jwtMocks.JWTValidator)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success - Get Bookmarks Pagination",
			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=10", nil)
				req.Header.Set("Authorization", testTokenHeader)

				recorder := httptest.NewRecorder()
				app.ServeHTTP(recorder, req)
				return recorder
			},
			setupDB: func() *gorm.DB {
				db := fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
				// Seed a record under the authorized user id
				db.Create(&model.Bookmark{
					Base:        model.Base{ID: targetBookmarkUUID},
					Description: "GitHub",
					URL:         "https://github.com",
					Code:        "git12345",
					UserID:      uuid.MustParse(mockAuthUserID),
				})
				return db
			},
			setupMockJwt: func(mockVal *jwtMocks.JWTValidator) {
				// Import your JWT package if needed (e.g., "github.com/golang-jwt/jwt/v5")
				// active claims map matching what your middleware reads
				mockClaims := jwtutils.GetMapClaim(mockAuthUserID, testUsername)

				mockVal.On("ValidateJWT", "test-valid-jwt-token").Return(mockClaims, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `"count":3`,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// 1. Reset database fixtures
			db := tc.setupDB()
			redisMocks := pkgRedis.InitMockRedis(t)

			// 2. Build mock validator for privateBase protection layers
			mockJwtValidator := jwtMocks.NewJWTValidator(t)
			tc.setupMockJwt(mockJwtValidator)

			// 3. Mount engine with mock dependencies injected
			testRouter := api.NewEngine(&api.EngineOpts{
				Engine: gin.New(),
				Cfg:    &config.Config{},
				Redis:  redisMocks,
				SqlDB:  db,
				JwtGen: nil, // Setup not required for standard resource routes
				JwtVal: mockJwtValidator,
			})

			// 4. Fire the test execution
			record := tc.setupTestHTTP(testRouter)

			// 5. Run validation assertions
			assert.Equal(t, tc.expectedStatusCode, record.Code)
			if tc.expectedResponseBody != "" {
				assert.Contains(t, record.Body.String(), tc.expectedResponseBody)
			}
		})
	}
}

func TestEndpoint_Bookmark_Edit(t *testing.T) {
	t.Parallel()

	// Seed tracking identifiers
	targetBookmarkUUID := uuid.New()

	testCases := []struct {
		name                 string
		setupTestHTTP        func(api api.Engine) *httptest.ResponseRecorder
		setupDB              func() *gorm.DB
		setupMockJwt         func(mockVal *jwtMocks.JWTValidator)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success - Edit Bookmark",
			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(
					http.MethodPut,
					"/v1/bookmarks/"+targetBookmarkUUID.String(),
					bytes.NewBuffer([]byte(`{
						"description": "Updated GitHub Link",
						"url": "https://github.com/trending"
					}`)))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", testTokenHeader)

				recorder := httptest.NewRecorder()
				app.ServeHTTP(recorder, req)
				return recorder
			},
			setupDB: func() *gorm.DB {
				db := fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
				db.Create(&model.Bookmark{
					Base:        model.Base{ID: targetBookmarkUUID},
					Description: "GitHub",
					URL:         "https://github.com",
					Code:        "git12345",
					UserID:      uuid.MustParse(mockAuthUserID),
				})
				return db
			},
			setupMockJwt: func(mockVal *jwtMocks.JWTValidator) {
				// Import your JWT package if needed (e.g., "github.com/golang-jwt/jwt/v5")
				// active claims map matching what your middleware reads
				mockClaims := jwtutils.GetMapClaim(mockAuthUserID, testUsername)

				mockVal.On("ValidateJWT", "test-valid-jwt-token").Return(mockClaims, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "", // Add specific match string if your update handler responds with standard messages
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			redisMocks := pkgRedis.InitMockRedis(t)

			// 1. Reset database fixtures
			db := tc.setupDB()

			// 2. Build mock validator for privateBase protection layers
			mockJwtValidator := jwtMocks.NewJWTValidator(t)
			tc.setupMockJwt(mockJwtValidator)

			// 3. Mount engine with mock dependencies injected
			testRouter := api.NewEngine(&api.EngineOpts{
				Engine: gin.New(),
				Cfg:    &config.Config{},
				Redis:  redisMocks,
				SqlDB:  db,
				JwtGen: nil, // Setup not required for standard resource routes
				JwtVal: mockJwtValidator,
			})

			// 4. Fire the test execution
			record := tc.setupTestHTTP(testRouter)

			// 5. Run validation assertions
			assert.Equal(t, tc.expectedStatusCode, record.Code)
			if tc.expectedResponseBody != "" {
				assert.Contains(t, record.Body.String(), tc.expectedResponseBody)
			}
		})
	}
}

func TestEndpoint_Bookmark_Delete(t *testing.T) {
	t.Parallel()

	// Seed tracking identifiers
	targetBookmarkUUID := uuid.New()

	testCases := []struct {
		name                 string
		setupTestHTTP        func(api api.Engine) *httptest.ResponseRecorder
		setupDB              func() *gorm.DB
		setupMockJwt         func(mockVal *jwtMocks.JWTValidator)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success - Delete Bookmark",
			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(http.MethodDelete, "/v1/bookmarks/"+targetBookmarkUUID.String(), nil)
				req.Header.Set("Authorization", testTokenHeader)

				recorder := httptest.NewRecorder()
				app.ServeHTTP(recorder, req)
				return recorder
			},
			setupDB: func() *gorm.DB {
				db := fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
				db.Create(&model.Bookmark{
					Base:        model.Base{ID: targetBookmarkUUID},
					Description: "To Be Deleted",
					URL:         "https://temp.com",
					Code:        "tmp12345",
					UserID:      uuid.MustParse(mockAuthUserID),
				})
				return db
			},
			setupMockJwt: func(mockVal *jwtMocks.JWTValidator) {
				// Import your JWT package if needed (e.g., "github.com/golang-jwt/jwt/v5")
				// active claims map matching what your middleware reads
				mockClaims := jwtutils.GetMapClaim(mockAuthUserID, testUsername)

				mockVal.On("ValidateJWT", "test-valid-jwt-token").Return(mockClaims, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			redisMocks := pkgRedis.InitMockRedis(t)

			// 1. Reset database fixtures
			db := tc.setupDB()

			// 2. Build mock validator for privateBase protection layers
			mockJwtValidator := jwtMocks.NewJWTValidator(t)
			tc.setupMockJwt(mockJwtValidator)

			// 3. Mount engine with mock dependencies injected
			testRouter := api.NewEngine(&api.EngineOpts{
				Engine: gin.New(),
				Cfg:    &config.Config{},
				Redis:  redisMocks,
				SqlDB:  db,
				JwtGen: nil, // Setup not required for standard resource routes
				JwtVal: mockJwtValidator,
			})

			// 4. Fire the test execution
			record := tc.setupTestHTTP(testRouter)

			// 5. Run validation assertions
			assert.Equal(t, tc.expectedStatusCode, record.Code)
			if tc.expectedResponseBody != "" {
				assert.Contains(t, record.Body.String(), tc.expectedResponseBody)
			}
		})
	}
}

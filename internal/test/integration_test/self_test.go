package intergration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/api"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	jwtMocks "github.com/viettrung2103/bookmark-management/pkg/jwtutils/mocks"
	"gorm.io/gorm"
)

// TestEndpoint_SelfInfo tests the SelfInfo endpoint
func TestEndpoint_SelfInfo(t *testing.T) {
	t.Parallel()

	// Target user details sourced from your fixtures.UserCommonTestDB
	targetUserID := "87a3cb94-d2e8-422d-bb91-fc5215949eb8" // Jane Smith
	//invaliJWTToken := "123"

	testCases := []struct {
		name                 string
		setupTestHTTP        func(app api.Engine) *httptest.ResponseRecorder
		setupDB              func() *gorm.DB
		setupMockJwtVal      func(mockJwtVal *jwtMocks.JWTValidator)
		expectedStatusCode   int
		expectedResponseBody string
		verifyJSONResponse   func(t *testing.T, body string)
	}{
		{
			name: "success - authenticated user fetched profile",
			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(http.MethodGet, "/v1/self/info", nil)
				// Attach a dummy authorization token header so the middleware triggers
				req.Header.Set("Authorization", "Bearer valid-mock-token")

				recorder := httptest.NewRecorder()
				app.ServeHTTP(recorder, req)
				return recorder
			},
			setupDB: func() *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupMockJwtVal: func(mockJwtVal *jwtMocks.JWTValidator) {
				// The middleware calls ValidateJWT. We tell our mock to pretend the token
				// belongs to Jane Smith's ID.
				// Adjust "uid" to match the actual key string used inside your JWT claims.
				mockClaims := jwt.MapClaims{
					"uid":      targetUserID,
					"username": "janesmith_dev",
				}
				mockJwtVal.On("ValidateJWT", "valid-mock-token").Return(mockClaims, nil)
			},
			expectedStatusCode: http.StatusOK,
			verifyJSONResponse: func(t *testing.T, body string) {
				//var returnedUser model.User
				var response struct {
					Data model.User `json:"data"`
				}
				err := json.Unmarshal([]byte(body), &response)
				assert.NoError(t, err)
				//println(returnedUser)

				// Assert the handler output matched Jane Smith's information in the DB
				assert.Equal(t, targetUserID, response.Data.ID)
				assert.Equal(t, "Jane Smith", response.Data.DisplayName)
				assert.Equal(t, "janesmith_dev", response.Data.Username)
				assert.Equal(t, "jane.smith@example.com", response.Data.Email)
			},
		},
		{
			name: "failure - user not found in database",
			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(http.MethodGet, "/v1/self/info", nil)
				req.Header.Set("Authorization", `Bearer ghost-token`)

				recorder := httptest.NewRecorder()
				app.ServeHTTP(recorder, req)
				return recorder
			},
			setupDB: func() *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupMockJwtVal: func(mockJwtVal *jwtMocks.JWTValidator) {
				// Simulate token matching a user ID that does not exist in the database
				mockClaims := jwt.MapClaims{
					"uid": "non-existent-uuid-value",
				}
				mockJwtVal.On("ValidateJWT", "ghost-token").Return(mockClaims, nil)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"User Not Found"}`,
		},
		{
			name: "failure - unauthorized token invalid",
			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(http.MethodGet, "/v1/self/info", nil)
				req.Header.Set("Authorization", "Bearer invalid-token")

				recorder := httptest.NewRecorder()
				app.ServeHTTP(recorder, req)
				return recorder
			},
			setupDB: func() *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupMockJwtVal: func(mockJwtVal *jwtMocks.JWTValidator) {
				// Simulate a bad token causing your middleware validation to fail
				mockJwtVal.On("ValidateJWT", "invalid-token").Return(nil, assert.AnError)
			},
			// Status code depends entirely on how your middleware handles authentication failures.
			// Change http.StatusUnauthorized if your custom middleware uses a different status code.
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := tc.setupDB()

			// 1. Setup the JWT Validator Mock
			mockJwtVal := jwtMocks.NewJWTValidator(t)
			tc.setupMockJwtVal(mockJwtVal)

			// 2. Build the router engine injecting our mock validator
			testRouter := api.NewEngine(&api.EngineOpts{
				Engine: gin.New(),
				Cfg:    &config.Config{},
				//Redis:  nil,
				SqlDB: db,
				//JwtGen: nil,        // Not needed for fetching profiles
				JwtVal: mockJwtVal, // Injected mock to bypass the authentication barrier
			})

			// 3. Issue Request
			record := tc.setupTestHTTP(testRouter)

			// 4. Evaluate asserts
			assert.Equal(t, tc.expectedStatusCode, record.Code)

			if tc.expectedResponseBody != "" {
				assert.Contains(t, record.Body.String(), tc.expectedResponseBody)
			}
			if tc.verifyJSONResponse != nil {
				tc.verifyJSONResponse(t, record.Body.String())
			}
		})
	}
}

// TestEndpoint_EditSelfInfo tests the EditSelfInfo endpoint
func TestEndpoint_EditSelfInfo(t *testing.T) {
	t.Parallel()

	// Target user details sourced from your fixtures.UserCommonTestDB
	targetUserID := "87a3cb94-d2e8-422d-bb91-fc5215949eb8" // Jane Smith

	testCases := []struct {
		name                 string
		setupTestHTTP        func(app api.Engine) *httptest.ResponseRecorder
		setupDB              func() *gorm.DB
		setupMockJwtVal      func(mockJwtVal *jwtMocks.JWTValidator)
		expectedStatusCode   int
		expectedResponseBody string
		verifyDB             func(t *testing.T, db *gorm.DB) // Added to check actual DB state
	}{
		{
			name: "success - user info updated",
			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(http.MethodPut, "/v1/self/info",
					bytes.NewBuffer([]byte(`{
						"display_name": "Jane Smith Updated",
						"email": "jane.updated@example.com"
					}`)))
				req.Header.Set("Authorization", "Bearer valid-mock-token")
				req.Header.Set("Content-Type", "application/json") // Don't forget this for POST/PUT!

				recorder := httptest.NewRecorder()
				app.ServeHTTP(recorder, req)
				return recorder
			},
			setupDB: func() *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupMockJwtVal: func(mockJwtVal *jwtMocks.JWTValidator) {
				// The middleware calls ValidateJWT. We tell our mock to pretend the token
				// belongs to Jane Smith's ID.
				// Adjust "uid" to match the actual key string used inside your JWT claims.
				mockClaims := jwt.MapClaims{
					"uid":      targetUserID,
					"username": "janesmith_dev",
				}
				mockJwtVal.On("ValidateJWT", "valid-mock-token").Return(mockClaims, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `"Edit current user successfully!"`,
			verifyDB: func(t *testing.T, db *gorm.DB) {
				// Verify the database actually saved the new values
				var user model.User
				err := db.First(&user, "id = ?", targetUserID).Error

				assert.NoError(t, err)
				assert.Equal(t, "Jane Smith Updated", user.DisplayName)
				assert.Equal(t, "jane.updated@example.com", user.Email)
			},
		},
		{
			name: "failure - validation error (invalid email format)",
			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(http.MethodPut, "/v1/self/info",
					bytes.NewBuffer([]byte(`{
						"display_name": "Jane Smith Updated",
						"email": "not email"
					}`)))
				req.Header.Set("Authorization", "Bearer valid-mock-token")
				req.Header.Set("Content-Type", "application/json") // Don't forget this for POST/PUT!

				recorder := httptest.NewRecorder()
				app.ServeHTTP(recorder, req)
				return recorder
			},
			setupDB: func() *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupMockJwtVal: func(mockJwtVal *jwtMocks.JWTValidator) {
				// The middleware calls ValidateJWT. We tell our mock to pretend the token
				// belongs to Jane Smith's ID.
				// Adjust "uid" to match the actual key string used inside your JWT claims.
				mockClaims := jwt.MapClaims{
					"uid":      targetUserID,
					"username": "janesmith_dev",
				}
				mockJwtVal.On("ValidateJWT", "valid-mock-token").Return(mockClaims, nil)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "",
			verifyDB: func(t *testing.T, db *gorm.DB) {
				var user model.User
				db.First(&user, "id = ?", targetUserID)
				assert.Equal(t, "jane.smith@example.com", user.Email)
			},
		},
		{
			name: "failure - duplication error (new email already taken)",
			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(http.MethodPut, "/v1/self/info",
					bytes.NewBuffer([]byte(`{
						"display_name": "Jane Smith Updated",
						"email": "john.doe@example.com"
					}`)))
				req.Header.Set("Authorization", "Bearer valid-mock-token")
				req.Header.Set("Content-Type", "application/json") // Don't forget this for POST/PUT!

				recorder := httptest.NewRecorder()
				app.ServeHTTP(recorder, req)
				return recorder
			},
			setupDB: func() *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupMockJwtVal: func(mockJwtVal *jwtMocks.JWTValidator) {
				// The middleware calls ValidateJWT. We tell our mock to pretend the token
				// belongs to Jane Smith's ID.
				// Adjust "uid" to match the actual key string used inside your JWT claims.
				mockClaims := jwt.MapClaims{
					"uid":      targetUserID,
					"username": "janesmith_dev",
				}
				mockJwtVal.On("ValidateJWT", "valid-mock-token").Return(mockClaims, nil)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "",
			verifyDB: func(t *testing.T, db *gorm.DB) {
				var user model.User
				db.First(&user, "id = ?", targetUserID)
				assert.Equal(t, "jane.smith@example.com", user.Email)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := tc.setupDB()

			// 1. Setup the JWT Validator Mock
			mockJwtVal := jwtMocks.NewJWTValidator(t)
			tc.setupMockJwtVal(mockJwtVal)

			// 2. Build the router engine injecting our mock validator
			testRouter := api.NewEngine(&api.EngineOpts{
				Engine: gin.New(),
				Cfg:    &config.Config{},
				//Redis:  nil,
				SqlDB: db,
				//JwtGen: nil,        // Not needed for fetching profiles
				JwtVal: mockJwtVal, // Injected mock to bypass the authentication barrier
			})

			// 3. Issue Request
			record := tc.setupTestHTTP(testRouter)

			// 4. Evaluate asserts
			assert.Equal(t, tc.expectedStatusCode, record.Code)

			if tc.expectedResponseBody != "" {
				assert.Contains(t, record.Body.String(), tc.expectedResponseBody)
			}
			if tc.verifyDB != nil {
				tc.verifyDB(t, db)
			}
		})
	}
}

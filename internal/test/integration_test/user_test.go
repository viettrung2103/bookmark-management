package intergration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	//"github.com/stretchr/testify/mock"
	"github.com/viettrung2103/bookmark-management/internal/api"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	jwtMocks "github.com/viettrung2103/bookmark-management/pkg/jwtutils/mocks"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"
	"gorm.io/gorm"
)

// TestRegisterEndpoint tests the register endpoint
func TestRegisterEndpoint(t *testing.T) {
	t.Parallel()

	//// Standard test data
	//testUUID := "12345678-1234-1234-1234-123456789012"
	//
	//username := "testuser"
	//password := "plainpassword"
	//hashedPassword := "hashedpassword"
	//mockUser := &model.User{
	//	Base: model.Base{
	//		ID: uuid.MustParse(testUUID),
	//	},
	//	//ID:       "user-123",
	//	Username: username,
	//	Password: hashedPassword,
	//}
	//expectedToken := "fake-jwt-token"

	testCases := []struct {
		name                 string
		setupTestHTTP        func(api api.Engine) *httptest.ResponseRecorder
		setupDB              func() *gorm.DB
		expectedErrString    string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "normal case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(
					"POST",
					"/v1/users/register",
					bytes.NewBuffer([]byte(
						`{
					"display_name":"test",
"username":"test1",
"password":"test12345",
"email":"test1@mail.com"
}`)))
				req.Header.Set("Content-Type", "application/json")

				respondRecorder := httptest.NewRecorder()
				api.ServeHTTP(respondRecorder, req)
				return respondRecorder
			},
			setupDB: func() *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			expectedStatusCode:   http.StatusOK,
			expectedErrString:    "",
			expectedResponseBody: "Register an user successfully",
		},
		{
			name: "err case - unique username",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(
					"POST",
					"/v1/users/register",
					bytes.NewBuffer([]byte(
						`{
					"display_name":"test",
"username":"test",
"password":"test12345",
"email":"test1@mail.com"
}`)))
				req.Header.Set("Content-Type", "application/json")

				resRecorde := httptest.NewRecorder()
				api.ServeHTTP(resRecorde, req)
				return resRecorde
			},
			setupDB: func() *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedErrString:    "UNIQUE",
			expectedResponseBody: `{"message":"User or Email already exist"}`,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// generate test cache
			fixtures := tc.setupDB()

			testRouter := api.NewEngine(&api.EngineOpts{
				Engine: gin.New(),
				Cfg:    &config.Config{},
				Redis:  nil,
				SqlDB:  fixtures,
			})

			record := tc.setupTestHTTP(testRouter)

			assert.Equal(t, tc.expectedStatusCode, record.Code)
			assert.Contains(t, record.Body.String(), tc.expectedResponseBody)
		})
	}
}

// TestEngine_Login tests the login endpoint
func TestEngine_Login(t *testing.T) {
	t.Parallel()
	// Standard test data
	testUUID := "12345678-1234-1234-1234-123456789012"
	displayName := "testuser"
	username := "testuser"
	password := "plainpassword"
	hashedPassword := "hashedpassword"
	email := "testuser@mail.com"
	mockUser := &model.User{
		Base: model.Base{
			ID: uuid.MustParse(testUUID),
		},
		//ID:       "user-123",
		Username: username,
		Password: hashedPassword,
	}
	expectedToken := "fake-jwt-token"

	testCases := []struct {
		name                 string
		setupTestHTTP        func(api api.Engine) *httptest.ResponseRecorder
		setupDB              func() *gorm.DB
		setupMockJwt         func(mockJwt *jwtMocks.JWTGenerator)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "success - valid login",
			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(
					http.MethodPost,
					"/v1/users/login",
					bytes.NewBuffer([]byte(`{
						"username": "testuser",
						"password": "plainpassword"
					}`)))
				req.Header.Set("Content-Type", "application/json")

				recorder := httptest.NewRecorder()
				app.ServeHTTP(recorder, req)
				return recorder
			},
			setupDB: func() *gorm.DB {
				db := fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})

				// 1. Generate a real hash for a password we actually know
				hasher := stringutils.NewPasswordHasher()
				knownHash := hasher.Hashing(password)

				// 2. Insert this specific user into the test database
				db.Create(&model.User{
					//ID:          "integration-test-id",
					Base: model.Base{
						ID: uuid.MustParse(testUUID),
					},
					DisplayName: displayName,
					Username:    username,
					Password:    knownHash,
					Email:       email,
				})

				return db
			},
			setupMockJwt: func(mockJwt *jwtMocks.JWTGenerator) {
				// Simulate successful token generation
				//testTokenClaim := jwtutils.GetMapClaim(mockUser.ID.String(), mockUser.Username)
				mockJwt.On("GenerateJWT", mock.MatchedBy(func(c jwt.MapClaims) bool {
					return c["uid"] == mockUser.ID.String() &&
						c["username"] == mockUser.Username
				})).Return(expectedToken, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"data":"fake-jwt-token","message":"Logged in successfully!"}`,
		},
		{
			name: "failure - wrong password",
			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req, _ := http.NewRequest(
					http.MethodPost,
					"/v1/users/login",
					bytes.NewBuffer([]byte(`{
						"username": "integration_user",
						"password": "wrongpassword999"
					}`)))
				req.Header.Set("Content-Type", "application/json")

				recorder := httptest.NewRecorder()
				app.ServeHTTP(recorder, req)
				return recorder
			},
			setupDB: func() *gorm.DB {
				db := fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})

				// Insert the same known user so the database finds them,
				// but the password check will fail in the service layer.
				hasher := stringutils.NewPasswordHasher()
				knownHash := hasher.Hashing("my_known_password123")

				db.Create(&model.User{
					//ID:          "integration-test-id",
					Base: model.Base{
						ID: uuid.MustParse(testUUID),
					},
					DisplayName: "Integration User",
					Username:    "integration_user",
					Password:    knownHash,
					Email:       "integration@example.com",
				})

				return db
			},
			setupMockJwt: func(mockJwt *jwtMocks.JWTGenerator) {
				// We don't expect the JWT generator to be called
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"Invalid username or password"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// 1. Generate test database via fixture
			db := tc.setupDB()

			// 2. Setup JWT mock
			mockJwt := jwtMocks.NewJWTGenerator(t)
			tc.setupMockJwt(mockJwt)

			// 3. Initialize engine with the injected test dependencies
			testRouter := api.NewEngine(&api.EngineOpts{
				Engine: gin.New(),
				Cfg:    &config.Config{},
				Redis:  nil,
				SqlDB:  db,
				JwtGen: mockJwt, // Inject the mock so token generation doesn't panic
				JwtVal: nil,     // Not needed for Login endpoint
			})

			// 4. Run the HTTP request
			record := tc.setupTestHTTP(testRouter)

			// 5. Assert outcomes
			assert.Equal(t, tc.expectedStatusCode, record.Code)
			assert.Contains(t, record.Body.String(), tc.expectedResponseBody)
		})
	}
}

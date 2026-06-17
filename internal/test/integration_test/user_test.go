package intergration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/api"
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	"gorm.io/gorm"
)

// TestRegisterEndpoint tests the register endpoint
func TestRegisterEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder
		setupDB       func() *gorm.DB

		expectedErrString string

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
			name: "err case - unique usernam",
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
			expectedStatusCode:   http.StatusInternalServerError,
			expectedErrString:    "UNIQUE",
			expectedResponseBody: `{"error":"Cannot create user, please try again"}`,
		},
		//		{
		//			name: "err case - unique email",
		//			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
		//				req, _ := http.NewRequest(
		//					"POST",
		//					"/v1/users/register",
		//					bytes.NewBuffer([]byte(
		//						`{
		//					"display_name":"test",
		//"username":"test1",
		//"password":"test12345",
		//"email":"test@mail.com"
		//}`)))
		//				req.Header.Set("Content-Type", "application/json")
		//
		//				resRecorde := httptest.NewRecorder()
		//				api.ServeHTTP(resRecorde, req)
		//				return resRecorde
		//			},
		//			setupDB: func() *gorm.DB {
		//				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
		//			},
		//			expectedStatusCode:   http.StatusBadRequest,
		//			expectedResponseBody: `{"message":"Field users.email already exist"}`,
		//		},
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

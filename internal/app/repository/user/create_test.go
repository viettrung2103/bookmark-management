package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	"gorm.io/gorm"
)

func TestUserRepo_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB           func(t *testing.T) *gorm.DB
		expectedErrString string
		inputUser         *model.User

		verifyFunc func(db *gorm.DB)
	}{
		{
			name: "normal case",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},

			inputUser: &model.User{
				ID:          "87a3cb94-d2e8-422d-bb91-fc5215949eb12",
				DisplayName: "Jane Smith43",
				Username:    "janesmith_dev123",
				Password:    "$2a$12$K3vX5YwOmP2bZ9rQ3nU7Xe8YvMw9T6uC9iK2oP1lRmSzTxVuWxYz.", // "hashed_password_2"
				Email:       "jane.smit23h@example.com",
			},
			expectedErrString: "",
			verifyFunc: func(db *gorm.DB) {
				user := &model.User{}
				err := db.First(user, "username = ? ", "janesmith_dev123").Error
				assert.NoError(t, err)
				assert.Equal(t, user.ID, "87a3cb94-d2e8-422d-bb91-fc5215949eb12")
				assert.Equal(t, user.DisplayName, "Jane Smith43")
				assert.Equal(t, user.Username, "janesmith_dev123")
				assert.Equal(t, user.Password, "$2a$12$K3vX5YwOmP2bZ9rQ3nU7Xe8YvMw9T6uC9iK2oP1lRmSzTxVuWxYz.")
				assert.Equal(t, user.Email, "jane.smit23h@example.com")
			},
		},
		{
			name: "same user error case",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},

			inputUser: &model.User{
				ID:          "87a3cb94-d2e8-422d-bb91-fc5215949eb8",
				DisplayName: "Jane Smith",
				Username:    "janesmith_dev",
				Password:    "$2a$12$K3vX5YwOmP2bZ9rQ3nU7Xe8YvMw9T6uC9iK2oP1lRmSzTxVuWxYz.", // "hashed_password_2"
				Email:       "jane.smith@example.com",
			},
			expectedErrString: "UNIQUE",
			verifyFunc: func(db *gorm.DB) {

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)

			repo := NewRepository(db)

			err := repo.CreateUser(ctx, tc.inputUser)
			//if tc.expectedErr {
			//	assert.NotNil(t, err)
			//} else {
			//	assert.NoError(t, err)
			//}
			if tc.expectedErrString == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.expectedErrString)
			}
			tc.verifyFunc(db)

		})
	}
}

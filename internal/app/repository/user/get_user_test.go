package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

func TestUserRepo_GetUserByUserId(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupDB       func(t *testing.T) *gorm.DB
		inputUserId   string
		expectedError error
		verifyUser    func(t *testing.T, user *model.User)
	}{
		{
			name: "success - user found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			// Using the ID from your previous fixture example
			inputUserId:   "133b3b42-70b9-456c-82e7-bf1b570e6c51",
			expectedError: nil,
			verifyUser: func(t *testing.T, user *model.User) {
				assert.NotNil(t, user)
				assert.Equal(t, "133b3b42-70b9-456c-82e7-bf1b570e6c51", user.ID.String())
			},
		},
		{
			name: "failure - user not found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUserId:   "invalid-id-999",
			expectedError: dbutils.ErrRecordNotFound,
			verifyUser: func(t *testing.T, user *model.User) {
				assert.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewRepository(db)

			user, err := repo.GetUserByUserId(ctx, tc.inputUserId)

			if tc.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.expectedError)
			}

			tc.verifyUser(t, user)
		})
	}
}

func TestUserRepo_GetUserByUsername(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupDB       func(t *testing.T) *gorm.DB
		inputUsername string
		expectedError error
		verifyUser    func(t *testing.T, user *model.User)
	}{
		{
			name: "success - user found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			// Using the username from your previous fixture example
			inputUsername: "janesmith_dev",
			expectedError: nil,
			verifyUser: func(t *testing.T, user *model.User) {
				assert.NotNil(t, user)
				assert.Equal(t, "janesmith_dev", user.Username)
			},
		},
		{
			name: "failure - user not found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUsername: "non_existent_username",
			// Assuming dbutils.CatchDBError translates gorm.ErrRecordNotFound to dbutils.ErrRecordNotFound
			expectedError: dbutils.ErrRecordNotFound,
			verifyUser: func(t *testing.T, user *model.User) {
				assert.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewRepository(db)

			user, err := repo.GetUserByUsername(ctx, tc.inputUsername)

			if tc.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.expectedError)
			}

			tc.verifyUser(t, user)
		})
	}
}

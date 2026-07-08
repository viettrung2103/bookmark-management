package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

// TestUserRepo_EditUserByID tests the EditUserByID method
func TestUserRepo_EditUserByID(t *testing.T) {
	t.Parallel()
	foundUUID := "133b3b42-70b9-456c-82e7-bf1b570e6c51"
	invalidUUID := "00000000-0000-0000-0000-000000000000"

	testCases := []struct {
		name             string
		setupDB          func(t *testing.T) *gorm.DB
		inputUserID      string
		inputDisplayName string
		inputEmail       string
		expectedError    error
		verifyDB         func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "success - user updated",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			// Using Jane Smith's ID from your fixture
			inputUserID:      foundUUID,
			inputDisplayName: "Jane Smith Updated",
			inputEmail:       "jane.updated@example.com",
			expectedError:    nil,
			verifyDB: func(t *testing.T, db *gorm.DB) {
				// Query the DB directly to verify the update persisted
				var updatedUser model.User
				err := db.First(&updatedUser, "id = ?", foundUUID).Error

				assert.NoError(t, err)
				assert.Equal(t, "Jane Smith Updated", updatedUser.DisplayName)
				assert.Equal(t, "jane.updated@example.com", updatedUser.Email)
			},
		},
		{
			name: "failure - user not found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUserID:      invalidUUID,
			inputDisplayName: "Ghost User",
			inputEmail:       "ghost@example.com",
			expectedError:    dbutils.ErrRecordNotFound, // Triggered by your RowsAffected == 0 check
			verifyDB: func(t *testing.T, db *gorm.DB) {
				// Nothing to verify here since no record was updated
			},
		},
		{
			name: "failure - duplication error (email already taken)",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			// Using Jane Smith's ID again
			inputUserID:      "87a3cb94-d2e8-422d-bb91-fc5215949eb8",
			inputDisplayName: "Jane Smith",
			// Attempting to change to John Doe's email (which exists in the fixture)
			inputEmail:    "john.doe@example.com",
			expectedError: dbutils.ErrDuplication, // Triggered by GORM error mapped via CatchDBError
			verifyDB: func(t *testing.T, db *gorm.DB) {
				// Verify the rollback: Jane's email should remain unchanged
				var user model.User
				err := db.First(&user, "id = ?", "87a3cb94-d2e8-422d-bb91-fc5215949eb8").Error

				assert.NoError(t, err)
				assert.Equal(t, "jane.smith@example.com", user.Email) // Still the old email
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewRepository(db)

			// Execute the method
			err := repo.EditUserByID(ctx, tc.inputUserID, tc.inputDisplayName, tc.inputEmail)

			// Assert the error
			if tc.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.expectedError)
			}

			// Verify the state of the database post-execution
			tc.verifyDB(t, db)
		})
	}
}

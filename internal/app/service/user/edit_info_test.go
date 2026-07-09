package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	repoMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/user/mocks"

	//repoMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/mocks"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	hashingMocks "github.com/viettrung2103/bookmark-management/pkg/stringutils/mocks"
)

//var testUUID = "12345678-1234-1234-1234-123456789012"

func TestUserService_EditInfoByID(t *testing.T) {
	t.Parallel()

	// Standard inputs to use across all test cases
	userId := "user-123"
	inputDisplayName := "Updated Name"
	inputEmail := "updated@example.com"

	testCases := []struct {
		name          string
		setupMocks    func(repo *repoMocks.UserRepository)
		expectedError error
	}{
		{
			name: "success - info updated successfully",
			setupMocks: func(repo *repoMocks.UserRepository) {
				// Expect the exact parameters to be passed to the repo
				repo.On("EditUserByID", mock.Anything, userId, inputDisplayName, inputEmail).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "failure - user not found",
			setupMocks: func(repo *repoMocks.UserRepository) {
				repo.On("EditUserByID", mock.Anything, userId, inputDisplayName, inputEmail).Return(dbutils.ErrRecordNotFound)
			},
			expectedError: dbutils.ErrRecordNotFound,
		},
		{
			name: "failure - duplication error (email already taken)",
			setupMocks: func(repo *repoMocks.UserRepository) {
				repo.On("EditUserByID", mock.Anything, userId, inputDisplayName, inputEmail).Return(dbutils.ErrDuplication)
			},
			expectedError: dbutils.ErrDuplication,
		},
		{
			name: "failure - generic database error",
			setupMocks: func(repo *repoMocks.UserRepository) {
				repo.On("EditUserByID", mock.Anything, userId, inputDisplayName, inputEmail).Return(assert.AnError)
			},
			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Initialize mocks
			mockRepo := repoMocks.NewUserRepository(t)

			// We initialize these because NewService requires them,
			// but we don't need to set up any `.On()` expectations because
			// this specific method doesn't use the hasher or JWT generator.
			mockHasher := hashingMocks.NewPasswordHashing(t)
			//mockJwtGen := mocks.NewJWTGenerator(t)

			// Setup the specific repo mock for this scenario
			tc.setupMocks(mockRepo)

			// Initialize service
			opts := &UserServiceOpts{
				UserRepo:        mockRepo,
				PasswordHashing: mockHasher,
				//JwtGenerator:    mockJwtGen,
			}
			svc := NewService(opts)

			// Execute the method
			err := svc.EditInfoByID(ctx, userId, inputDisplayName, inputEmail)

			// Assert the results
			if tc.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.expectedError)
			}
		})
	}
}

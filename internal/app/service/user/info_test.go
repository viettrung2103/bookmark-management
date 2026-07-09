package user

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	//"github.com/viettrung2103/bookmark-management/internal/app/repository/user"
	repoMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/user/mocks"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	//hashingMocks"github.com/viettrung2103/bookmark-management/pkg/stringutils/mocks"
)

func TestUserService_SelfInfo(t *testing.T) {
	t.Parallel()

	testUUID := "12345678-1234-1234-1234-123456789012"

	// Create a mock user to return in the success case
	mockUser := &model.User{
		Base: model.Base{
			ID: uuid.MustParse(testUUID),
		},
		//ID:          userId,
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
	}

	testCases := []struct {
		name           string
		setupMocksRepo func(repo *repoMocks.UserRepository)
		expectedUser   *model.User
		expectedError  error
	}{
		{
			name: "success - user found",
			setupMocksRepo: func(repo *repoMocks.UserRepository) {
				repo.On("GetUserByUserId", mock.Anything, testUUID).Return(mockUser, nil)
			},
			expectedUser:  mockUser,
			expectedError: nil,
		},
		{
			name: "failure - user not found",
			setupMocksRepo: func(repo *repoMocks.UserRepository) {
				repo.On("GetUserByUserId", mock.Anything, testUUID).Return((*model.User)(nil), dbutils.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: dbutils.ErrRecordNotFound,
		},
		{
			name: "failure - generic database error",
			setupMocksRepo: func(repo *repoMocks.UserRepository) {
				repo.On("GetUserByUserId", mock.Anything, testUUID).Return((*model.User)(nil), assert.AnError)
			},
			expectedUser:  nil,
			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Initialize mocks
			mockRepo := repoMocks.NewUserRepository(t)
			//mockHasher := mocks.NewPasswordHashing(t) // Not used in this method
			//mockJwtGen := mocks.NewJWTGenerator(t)    // Not used in this method

			// Setup the specific repo mock for this scenario
			tc.setupMocksRepo(mockRepo)

			// Initialize service
			opts := &UserServiceOpts{
				UserRepo: mockRepo,
				//PasswordHashing: mockHasher,
				//JwtGenerator:    mockJwtGen,
			}
			svc := NewService(opts)

			// Execute the method
			user, err := svc.SelfInfo(ctx, testUUID)

			// Assert the results
			if tc.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
			} else {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, user)
			}
		})
	}
}

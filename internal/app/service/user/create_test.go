package user

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	repoMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/user/mocks"
	//serviceMocks"github.com/viettrung2103/bookmark-management/internal/app/service/mocks"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	hashingMocks "github.com/viettrung2103/bookmark-management/pkg/stringutils/mocks"
)

var testUUID = "12345678-1234-1234-1234-123456789012"

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()

	// Define standard inputs for our test cases
	displayName := "Jane Smith"
	username := "janesmith_dev"
	password := "supersecret123"
	email := "jane@example.com"
	hashedPassword := "hashed_supersecret123"

	// The exact user struct the service should build and pass to the repository
	expectedUserToSave := &model.User{
		Username:    username,
		Password:    hashedPassword,
		Email:       email,
		DisplayName: displayName,
	}

	testCases := []struct {
		name          string
		setupMocks    func(ctx context.Context, repo *repoMocks.UserRepository, hasher *hashingMocks.PasswordHashing)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "success - user created successfully",
			setupMocks: func(ctx context.Context, repo *repoMocks.UserRepository, hasher *hashingMocks.PasswordHashing) {
				// 1. Mock the hashing utility
				hasher.On("Hashing", password).Return(hashedPassword)

				// 2. Mock the repository saving the exact user struct
				mockSavedUser := &model.User{
					Base: model.Base{
						ID: uuid.MustParse(testUUID),
					},
					//ID:          "user-123",
					Username:    username,
					Password:    hashedPassword,
					Email:       email,
					DisplayName: displayName,
				}
				repo.On("CreateUser", ctx, expectedUserToSave).Return(mockSavedUser, nil)
			},
			expectedUser: &model.User{
				Base: model.Base{
					ID: uuid.MustParse(testUUID),
				},
				//ID:          "user-123",
				Username:    username,
				Password:    hashedPassword,
				Email:       email,
				DisplayName: displayName,
			},
			expectedError: nil,
		},
		{
			name: "failure - repository returns duplication error",
			setupMocks: func(ctx context.Context, repo *repoMocks.UserRepository, hasher *hashingMocks.PasswordHashing) {
				// 1. Hashing still succeeds
				hasher.On("Hashing", password).Return(hashedPassword)

				// 2. Repository simulates a unique constraint failure (e.g., email already exists)
				repo.On("CreateUser", ctx, expectedUserToSave).Return((*model.User)(nil), dbutils.ErrDuplication)
			},
			expectedUser:  nil,
			expectedError: dbutils.ErrDuplication,
		},
		{
			name: "failure - repository returns generic error",
			setupMocks: func(ctx context.Context, repo *repoMocks.UserRepository, hasher *hashingMocks.PasswordHashing) {
				hasher.On("Hashing", password).Return(hashedPassword)

				// Simulate a database crash
				repo.On("CreateUser", ctx, expectedUserToSave).Return((*model.User)(nil), assert.AnError)
			},
			expectedUser:  nil,
			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			// 1. Initialize mocks
			mockRepo := repoMocks.NewUserRepository(t)
			mockHasher := hashingMocks.NewPasswordHashing(t)
			//mockJwtGen := mocks.NewJWTGenerator(t) // Not used in this method, but required by UserServiceOpts

			// 2. Run test case specific mock setup
			tc.setupMocks(ctx, mockRepo, mockHasher)

			// 3. Initialize service using your opts struct
			opts := &UserServiceOpts{
				UserRepo:        mockRepo,
				PasswordHashing: mockHasher,
			}
			svc := NewService(opts)

			// 4. Execute the method
			resultUser, err := svc.CreateUser(ctx, displayName, username, password, email)

			// 5. Assert the results
			if tc.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, resultUser)
			} else {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, resultUser)
			}
		})
	}
}

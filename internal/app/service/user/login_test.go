package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	//repoMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/mocks"
	userMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/user/mocks"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	jwtMocks "github.com/viettrung2103/bookmark-management/pkg/jwtutils/mocks"
	hashingMocks "github.com/viettrung2103/bookmark-management/pkg/stringutils/mocks"
)

func TestUserService_Login(t *testing.T) {
	t.Parallel()

	// Standard test data
	testUUID := "12345678-1234-1234-1234-123456789012"

	username := "testuser"
	password := "plainpassword"
	hashedPassword := "hashedpassword"
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
		name          string
		setupMocks    func(repo *userMocks.UserRepository, hasher *hashingMocks.PasswordHashing, jwtGen *jwtMocks.JWTGenerator)
		expectedToken string
		expectedError error
	}{
		{
			name: "success - login successful",
			setupMocks: func(repo *userMocks.UserRepository, hasher *hashingMocks.PasswordHashing, jwtGen *jwtMocks.JWTGenerator) {
				// 1. Repo finds the user
				repo.On("GetUserByUsername", mock.Anything, username).Return(mockUser, nil)

				// 2. Password matches
				hasher.On("CompareHashedPassword", hashedPassword, password).Return(true)

				// 3. Token is generated
				jwtGen.On("GenerateJWT", mock.Anything).Return(expectedToken, nil)
			},
			expectedToken: expectedToken,
			expectedError: nil,
		},
		{
			name: "failure - user not found",
			setupMocks: func(repo *userMocks.UserRepository, hasher *hashingMocks.PasswordHashing, jwtGen *jwtMocks.JWTGenerator) {
				repo.On("GetUserByUsername", mock.Anything, username).Return((*model.User)(nil), dbutils.ErrRecordNotFound)
				// We don't need to mock hasher or jwtGen because the function returns early
			},
			expectedToken: "",
			expectedError: ErrInvalidCreditials, // Should bubble up the exact error from the DB, but for security reason, display same error
		},
		{
			name: "failure - wrong password",
			setupMocks: func(repo *userMocks.UserRepository, hasher *hashingMocks.PasswordHashing, jwtGen *jwtMocks.JWTGenerator) {
				// Repo finds the user
				repo.On("GetUserByUsername", mock.Anything, username).Return(mockUser, nil)

				// Password check fails (returns false)
				hasher.On("CompareHashedPassword", hashedPassword, password).Return(false)
			},
			expectedToken: "",
			expectedError: ErrInvalidCreditials,
		},
		{
			name: "failure - jwt generation fails",
			setupMocks: func(repo *userMocks.UserRepository, hasher *hashingMocks.PasswordHashing, jwtGen *jwtMocks.JWTGenerator) {
				repo.On("GetUserByUsername", mock.Anything, username).Return(mockUser, nil)
				hasher.On("CompareHashedPassword", hashedPassword, password).Return(true)

				// Token generator fails
				jwtGen.On("GenerateJWT", mock.Anything).Return("", errors.New("jwt signing error"))
			},
			expectedToken: "",
			expectedError: errors.New("jwt signing error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// 1. Initialize all mocks
			mockRepo := userMocks.NewUserRepository(t)
			mockHasher := hashingMocks.NewPasswordHashing(t)
			mockJwtGen := jwtMocks.NewJWTGenerator(t)

			// 2. Setup mock expectations
			tc.setupMocks(mockRepo, mockHasher, mockJwtGen)

			// 3. Initialize the service
			opts := &UserServiceOpts{
				UserRepo:        mockRepo,
				PasswordHashing: mockHasher,
				JwtGenerator:    mockJwtGen,
			}
			svc := NewService(opts)

			// 4. Call the method
			token, err := svc.Login(ctx, username, password)

			// 5. Assert the outcomes
			if tc.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedToken, token)
			} else {
				// Using EqualError for generic errors, ErrorIs would also work for sentinels
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Empty(t, token)
			}
		})
	}
}

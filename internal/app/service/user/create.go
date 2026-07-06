package user

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
)

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, displayName, username, password, email string) (*model.User, error) {
	hashedPwd := s.passwordHashing.Hashing(password)

	user := &model.User{
		Username:    username,
		Password:    hashedPwd,
		Email:       email,
		DisplayName: displayName,
	}

	newUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

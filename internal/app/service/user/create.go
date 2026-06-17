package user

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"
)

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, displayName, username, password, email string) (*model.User, error) {
	hashedPwd := stringutils.Hastring(password)

	user := &model.User{
		Username:    username,
		Password:    hashedPwd,
		Email:       email,
		DisplayName: displayName,
	}

	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

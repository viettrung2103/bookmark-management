package user

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/app/repository/user"
)

// Service interface for user service
type Service interface {
	CreateUser(ctx context.Context, displayName, username, password, email string) (*model.User, error)
}

type userService struct {
	repo user.Repository
}

// NewService creates a new user service
func NewService(repo user.Repository) Service {
	return &userService{repo: repo}
}

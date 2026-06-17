package user

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/app/repository/user"
)

type Service interface {
	CreateUser(ctx context.Context, displayName, username, password, email string) (*model.User, error)
}

type userService struct {
	repo user.Repository
}

func NewService(repo user.Repository) Service {
	return &userService{repo: repo}
}

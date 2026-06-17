package user

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"gorm.io/gorm"
)

// Repository interface
type Repository interface {
	CreateUser(ctx context.Context, user *model.User) error
}

type userRepo struct {
	db *gorm.DB
}

// NewRepository returns a new repository
func NewRepository(db *gorm.DB) Repository {
	return &userRepo{db: db}
}

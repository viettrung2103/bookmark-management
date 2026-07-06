package user

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
)

// CreateUser creates a new user
func (r *userRepo) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

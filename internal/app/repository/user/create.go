package user

import (
	"context"

	"github.com/viettrung2103/bookmark-lib/pkg/dbutils"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
)

// CreateUser creates a new user
func (r *userRepository) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).Create(newUser).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return newUser, nil
}

package user

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
)

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	err := r.db.WithContext(ctx).Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return user, nil
}

func (r *userRepository) GetUserByUserId(ctx context.Context, userId string) (*model.User, error) {
	user := &model.User{}
	err := r.db.WithContext(ctx).Where("id = ?", userId).First(user).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return user, nil
}

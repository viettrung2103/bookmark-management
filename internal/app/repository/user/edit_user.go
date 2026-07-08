package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/common"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

// EditUserByID edits a user by ID
func (r *userRepository) EditUserByID(ctx context.Context, userID string, inputDisplayName string, inputEmail string) error {

	parsedUserId, err := uuid.Parse(userID)
	common.HandleError(err)

	updateUser := model.User{
		DisplayName: inputDisplayName,
		Email:       inputEmail,
	}

	result := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id=?", parsedUserId).
		Select("DisplayName", "Email").
		Updates(&updateUser)

	if result.Error != nil {
		return dbutils.CatchDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return dbutils.CatchDBError(gorm.ErrRecordNotFound)
	}
	return nil
}

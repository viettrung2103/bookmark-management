package bookmark

import (
	"context"

	"github.com/google/uuid"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/common"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
)

// DeleteBookmarkByID delete bookmark of id of current user
func (r *bookmarkRepository) DeleteBookmarkByID(ctx context.Context, userID, bookmarkID string) error {

	parsedUserID, err := uuid.Parse(userID)
	common.HandleError(err)

	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	common.HandleError(err)

	// it return an transaction
	result := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("id=? AND user_id=?", parsedBookmarkID, parsedUserID).
		Delete(&model.Bookmark{})

	if result.Error != nil {
		return dbutils.CatchDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return dbutils.CatchDBError(dbutils.ErrRecordNotFound)
	}
	return nil
}

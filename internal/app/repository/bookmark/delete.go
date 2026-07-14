package bookmark

import (
	"context"

	"github.com/google/uuid"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/common"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
)

func (r *bookmarkRepository) DeleteBookmarkByID(ctx context.Context, userID, bookmarkID string) error {

	parsedUserID, err := uuid.Parse(userID)
	common.HandleError(err)

	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	common.HandleError(err)

	println("is problem in repo")

	// it return an transaction
	result := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("id=? AND user_id=?", parsedBookmarkID, parsedUserID).
		Delete(&model.Bookmark{})

	println("is problem in db")

	if result.Error != nil {
		return dbutils.CatchDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return dbutils.CatchDBError(dbutils.ErrRecordNotFound)
	}
	return nil
}

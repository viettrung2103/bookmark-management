package bookmark

import (
	"context"

	"github.com/google/uuid"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/common"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
)

// EditBookmarkByID edit bookmark of id of current user
func (r *bookmarkRepository) EditBookmarkByID(ctx context.Context, userID, bookmarkID, newDescription, newURL string) error {
	parsedUserID, err := uuid.Parse(userID)
	common.HandleError(err)

	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	common.HandleError(err)

	updatedBookmark := model.Bookmark{
		Description: newDescription,
		URL:         newURL,
	}

	result := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("id = ? AND user_id = ?", parsedBookmarkID, parsedUserID).
		Select("description", "url").
		Updates(&updatedBookmark)

	if result.Error != nil {
		return dbutils.CatchDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		//return dbutils.CatchDBError(gorm.ErrRecordNotFound)
		return dbutils.CatchDBError(dbutils.ErrRecordNotFound)
	}
	return nil
}

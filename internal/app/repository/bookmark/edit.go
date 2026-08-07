package bookmark

import (
	"context"

	"github.com/google/uuid"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/common"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
)

// EditBookmarkByID edit bookmark of id of current user
type ParsedBookmarkIDs struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (r *bookmarkRepository) EditBookmarkByID(ctx context.Context, userID, bookmarkID, newDescription, newURL string) error {

	parsedBookmarkID := getParsedUserIDAndParsedBookmarkID(bookmarkID, userID)

	updatedBookmark := model.Bookmark{
		Description: newDescription,
		URL:         newURL,
	}

	result := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("id = ? AND user_id = ?", parsedBookmarkID.ID, parsedBookmarkID.UserID).
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

func (r *bookmarkRepository) EditBookmarkCodeByID(ctx context.Context, bookmarkID, userID string, newCode string) error {
	parsedBookmarkID := getParsedUserIDAndParsedBookmarkID(bookmarkID, userID)

	updatedBookmark := model.Bookmark{
		Code: newCode,
	}

	result := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("id = ? AND user_id = ?", parsedBookmarkID.ID, parsedBookmarkID.UserID).
		Select("code").
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

func getParsedUserIDAndParsedBookmarkID(userID, bookmarkID string) ParsedBookmarkIDs {

	parsedUserID, err := uuid.Parse(userID)
	common.HandleError(err)

	parsedBookmarkID, err := uuid.Parse(bookmarkID)
	common.HandleError(err)
	return ParsedBookmarkIDs{
		ID:     parsedBookmarkID,
		UserID: parsedUserID,
	}
}

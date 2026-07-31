package bookmark

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
)

// GetBookmarks get bookmarks of current user
func (r *bookmarkRepository) GetBookmarks(ctx context.Context, userID string, offset, limit int) ([]*model.Bookmark, error) {
	// tao ra 1 array with fix size = limit
	bookmarks := make([]*model.Bookmark, 0, limit)
	err := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&bookmarks).
		Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return bookmarks, nil
}

// GetBookmarkCount get count of bookmarks of current user
func (r *bookmarkRepository) GetBookmarkCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("user_id = ?", userID).
		Count(&count).
		Error
	if err != nil {
		return 0, dbutils.CatchDBError(err)
	}

	return count, nil
}

func (r *bookmarkRepository) GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error) {
	println("get bookmark by code in bookmark repo")
	var selectedBookmark model.Bookmark
	// fetch bookmark from code
	err := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("code=?", code).
		First(&selectedBookmark).
		Error
	println("err", err)
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return &selectedBookmark, nil
}

package bookmark

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"gorm.io/gorm"
)

//go:generate mockery --name=Repository --filename=bookmark.go
type Repository interface {
	CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, offset, limit int) ([]*model.Bookmark, error)
	GetBookmarkCount(ctx context.Context, userID string) (int64, error)
	EditBookmarkByID(ctx context.Context, userID, bookmarkID, newDescription, newURL string) error
	DeleteBookmarkByID(ctx context.Context, userID, bookmarkID string) error
}

type bookmarkRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &bookmarkRepository{
		db: db,
	}
}

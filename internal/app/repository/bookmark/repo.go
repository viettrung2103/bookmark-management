package bookmark

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"
	"gorm.io/gorm"
)

// Repository interface for bookmark
//
//go:generate mockery --name=Repository --filename=bookmark.go
type Repository interface {
	CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, offset, limit int) ([]*model.Bookmark, error)
	GetBookmarkCount(ctx context.Context, userID string) (int64, error)
	GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error)
	EditBookmarkByID(ctx context.Context, userID, bookmarkID, newDescription, newURL string) error
	EditBookmarkCodeByID(ctx context.Context, userID, bookmarkID string, newCode string) error
	DeleteBookmarkByID(ctx context.Context, userID, bookmarkID string) error

	Transaction(ctx context.Context, fn func(txRepo Repository) error) error
}

type bookmarkRepository struct {
	db     *gorm.DB
	keygen stringutils.KeyGenerator
}

// NewRepository create new repository
func NewRepository(db *gorm.DB, keygen stringutils.KeyGenerator) Repository {
	return &bookmarkRepository{
		db:     db,
		keygen: keygen,
	}
}

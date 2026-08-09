package bookmark

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/app/repository/bookmark"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"
)

// BookmarkService interface for bookmark
//
//go:generate mockery --name=Service --filename=bookmark.go
type Service interface {
	AddBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarkResult, error)
	EditBookmarkByID(ctx context.Context, userID string, bookmarkID string, newDescription string, newURL string) error
	DeleteBookmarkByID(ctx context.Context, userID, bookmarkID string) error
}

// bookmarkService struct for bookmark
type bookmarkService struct {
	keygen       stringutils.KeyGenerator
	bookmarkRepo bookmark.Repository
}

// BookmarkServiceOpts contains dependencies for user service
type BookmarkServiceOpts struct {
	Keygen             stringutils.KeyGenerator
	BookmarkRepository bookmark.Repository
}

// NewService creates a new user service
func NewService(opts *BookmarkServiceOpts) Service {

	return &bookmarkService{
		keygen:       opts.Keygen,
		bookmarkRepo: opts.BookmarkRepository,
	}
}

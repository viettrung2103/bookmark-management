package bookmark

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/app/repository/bookmark"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"
)

//go:generate mockery --name=Service --filename=bookmark.go
type Service interface {
	CreateBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error)
	//Login(ctx context.Context, username, password string) (string, error)
	//SelfInfo(ctx context.Context, userId string) (*model.User, error)
	//EditInfoByID(ctx context.Context, userId string, inputDisplayName string, inputEmail string) error
}

type bookmarkService struct {
	keygen       stringutils.KeyGenerator
	bookmarkRepo bookmark.Repository
}

// UserServiceOpts contains dependencies for user service
type BookmarkServiceOpts struct {
	Keygen             stringutils.KeyGenerator
	BookmarkRepository bookmark.Repository
	//PasswordHashing stringutils.PasswordHashing
	//JwtGenerator    jwtutils.JWTGenerator
}

// NewService creates a new user service
func NewService(opts *BookmarkServiceOpts) Service {

	return &bookmarkService{
		keygen:       opts.Keygen,
		bookmarkRepo: opts.BookmarkRepository,
	}
}

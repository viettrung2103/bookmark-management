package bookmark

import (
	"context"
	"errors"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
)

var ErrNoOwnerShip = errors.New("bookmark does not belong to current user")

const shortenUrlKeyLength = 8

func (s *bookmarkService) AddBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	code := s.keygen.GenerateKey(shortenUrlKeyLength)

	newBookmark := &model.Bookmark{
		Description: description,
		URL:         url,
		Code:        code,
		UserID:      userID,
	}

	return s.bookmarkRepo.CreateBookmark(ctx, newBookmark)

}

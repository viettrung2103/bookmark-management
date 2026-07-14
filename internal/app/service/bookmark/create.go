package bookmark

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
)

const shortenUrlKeyLength = 8

// AddBookmark add new bookmark
func (s *bookmarkService) AddBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	code := s.keygen.GenerateKey(shortenUrlKeyLength)

	newBookmark := &model.Bookmark{
		Description: description,
		URL:         url,
		Code:        code,
		UserID:      fixtures.GetUUID(userID),
	}
	println("Add booking using repo")

	return s.bookmarkRepo.CreateBookmark(ctx, newBookmark)

}

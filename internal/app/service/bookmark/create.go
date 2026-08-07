package bookmark

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
)

const shortenUrlKeyLength = 8

// AddBookmark add new bookmark
func (s *bookmarkService) AddBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {

	var finalBookmark *model.Bookmark

	newBookmark := &model.Bookmark{
		Description: description,
		URL:         url,
		//Code:        code,
		UserID: fixtures.GetUUID(userID),
	}

	finalBookmark, err := s.bookmarkRepo.CreateBookmark(ctx, newBookmark)
	if err != nil {
		log.Err(err).Msg("failed to create bookmark")
		return nil, err
	}

	return finalBookmark, nil

}

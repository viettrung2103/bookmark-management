package bookmark

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
)

// GetBookmarkResult struct
type GetBookmarkResult struct {
	Bookmarks []*model.Bookmark `json:"bookmarks,omitempty"`
	Count     int64             `json:"count,omitempty"`
}

// GetBookmarks get lists of bookmarks of current user
func (s *bookmarkService) GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarkResult, error) {
	offset := (page - 1) * limit
	bookmarks, err := s.bookmarkRepo.GetBookmarks(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	count, err := s.bookmarkRepo.GetBookmarkCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &GetBookmarkResult{
		Bookmarks: bookmarks,
		Count:     count,
	}, nil

	//return bookmarks, nil
}

package bookmark

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	bookmarkRepository "github.com/viettrung2103/bookmark-management/internal/app/repository/bookmark"
	"github.com/viettrung2103/bookmark-management/internal/test/data/fixtures"
)

const shortenUrlKeyLength = 8

// AddBookmark add new bookmark
func (s *bookmarkService) AddBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {

	var finalBookmark *model.Bookmark

	// open transaction
	err := s.bookmarkRepo.Transaction(ctx, func(txRepo bookmarkRepository.Repository) error {
		code := s.keygen.GenerateKey(shortenUrlKeyLength)
		newBookmark := &model.Bookmark{
			Description: description,
			URL:         url,
			Code:        code,
			UserID:      fixtures.GetUUID(userID),
		}

		println("Add booking using repo")

		newBookmark, err := txRepo.CreateBookmark(ctx, newBookmark)
		println("new bookmark ", newBookmark)
		println("new bookmark & ", &newBookmark)
		if err != nil {
			log.Err(err).Msg("failed to create bookmark")
			return err // Phải kiểm tra lỗi ở đây để tránh panic!
		}

		// generate new code
		newCode := s.keygen.GenerateBase62Code(newBookmark.CodeInt)
		print("code ", newCode)
		newBookmark.Code = newCode

		// update new code to db
		err = txRepo.EditBookmarkCodeByID(ctx, newBookmark.ID.String(), newBookmark.UserID.String(), newCode)
		if err != nil {
			log.Err(err).Msg("failed to edit bookmark code by id")
			return err
		}

		// final bookmark and close transaction
		finalBookmark = newBookmark
		return nil
	})
	if err != nil {
		return nil, err
	}

	return finalBookmark, nil

}

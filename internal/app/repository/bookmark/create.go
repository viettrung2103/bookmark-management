package bookmark

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

// CreateBookmark create bookmark of current user
func (r *bookmarkRepository) CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) {
	//err := r.db.WithContext(ctx).Create(bookmark).Error
	//if err != nil {
	//	return nil, dbutils.CatchDBError(err)
	//}
	//
	////fmt.Printf("created bookmark: %+v\n", bookmark)
	//slog.Info("created bookmark", "bookmark", bookmark)
	//return bookmark, nil
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// insert initial record of bookmark
		if err := tx.Create(bookmark).Error; err != nil {
			return err
		}
		// create newcode
		newCode := r.keygen.GenerateBase62Code(bookmark.CodeInt)
		bookmark.Code = newCode

		// update new code to db
		if err := tx.Model(bookmark).Update("code", newCode).Error; err != nil {
			return err
		}

		// finish
		return nil
	})
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return bookmark, nil
}

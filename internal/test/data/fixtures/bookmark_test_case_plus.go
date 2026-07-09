package fixtures

import (
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"gorm.io/gorm"
)

type BookmarkPlusTestDB struct {
	BookmarkCommonTestDB
}

func (b *BookmarkPlusTestDB) GenerateData() error {

	err := b.BookmarkCommonTestDB.GenerateData()
	if err != nil {
		return err
	}

	bookmarks := []*model.Bookmark{
		{
			Base:        GetTestBase("133b3b42-70b9-456c-93d8-bf1b570e6c53"),
			Description: "bookmark1",
			URL:         "http://google.com",
			Code:        "123458",
			UserID:      GetUUID("133b3b42-70b9-456c-82e7-bf1b570e6c52"),
		},
		{
			Base:        GetTestBase("133b3b42-70b9-456c-93d8-bf1b570e6c54"),
			Description: "bookmark2",
			URL:         "http://google.com",
			Code:        "123459",
			UserID:      GetUUID("133b3b42-70b9-456c-82e7-bf1b570e6c52"),
		},
	}

	return b.db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(bookmarks, 10).Error
}

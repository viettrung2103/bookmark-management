package fixtures

import (
	//"github.com/google/uuid"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"gorm.io/gorm"
)

// UserCommonTestDB struct for user test data
type BookmarkCommonTestDB struct {
	UserCommonTestDB
}

// Migrate migrates the user table
func (b *BookmarkCommonTestDB) Migrate() error {

	// in case of we have a long migation
	//err := b.UserCommonTestDB.Migrate()
	//if err != nil {
	//	return err
	//}
	return b.db.AutoMigrate(&model.User{}, &model.Bookmark{})

	//return b.db.AutoMigrate(&model.Bookmark{})
}

// GenerateData generates test data for users and bookmark
func (b *BookmarkCommonTestDB) GenerateData() error {
	//db := b.db.Session(&gorm.Session{SkipHooks: true})

	err := b.UserCommonTestDB.GenerateData()
	if err != nil {
		return err
	}

	bookmarks := []*model.Bookmark{
		{
			Base:        GetTestBase("133b3b42-70b9-456c-93d8-bf1b570e6c55"),
			Description: "bookmark1",
			URL:         "http://google.com",
			Code:        "123456",
			UserID:      GetUUID("133b3b42-70b9-456c-82e7-bf1b570e6c51"),
		},
		{
			Base:        GetTestBase("133b3b42-70b9-456c-93d8-bf1b570e6c56"),
			Description: "bookmark2",
			URL:         "http://google.com",
			Code:        "123457",
			UserID:      GetUUID("133b3b42-70b9-456c-82e7-bf1b570e6c51"),
		},
	}

	return b.db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(bookmarks, 10).Error
}

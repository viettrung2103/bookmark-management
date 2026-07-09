package main

import (
	"fmt"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/common"
	"github.com/viettrung2103/bookmark-management/pkg/db"
)

func main() {
	sqlDB, err := db.NewClient("")
	common.HandleError(err)

	//db.MigrationPostgresDB(sqlDB, "up", 0)

	var bookmarks []*model.Bookmark

	err = sqlDB.Offset(15).Limit(10).Find(&bookmarks).Error
	println(err)
	common.HandleError(err)

	var total int64
	err = sqlDB.Model(&model.Bookmark{}).Count(&total).Error
	println(total)

	for _, bookmark := range bookmarks {
		fmt.Printf("%+v\n", bookmark)
	}

}

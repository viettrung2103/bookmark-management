package main

import (
	"github.com/viettrung2103/bookmark-management/pkg/common"
	"github.com/viettrung2103/bookmark-management/pkg/db"
)

func main() {
	sqlDB, err := db.NewClient("")
	common.HandleError(err)

	db.MigrationPostgresDB(sqlDB, "up", 0)
}

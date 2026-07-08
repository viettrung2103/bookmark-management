package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/viettrung2103/bookmark-management/pkg/common"
	"gorm.io/gorm"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migrationDir = "file://./migrations"
)

func MigrationPostgresDB(db *gorm.DB, mode string, step int) error {
	sqlDb, err := db.DB()
	common.HandleError(err)

	// create postgress driver
	driver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	common.HandleError(err)

	m, err := migrate.NewWithDatabaseInstance(migrationDir, db.Name(), driver)
	common.HandleError(err)

	return migrateSchema(m, mode, step)
}
func migrateSchema(m *migrate.Migrate, mode string, steps int) error {
	var migrationErr error

	// 2 mode, up and down
	// steps = 0, go to the latest up and down
	// steps > 0, migrate up/down that many steps
	switch mode {
	case "up":
		if steps > 0 {
			migrationErr = m.Steps(steps) // Move UP by 'steps'
		} else {
			migrationErr = m.Up() // Move UP completely
		}
	case "down":
		if steps > 0 {
			migrationErr = m.Steps(-steps) // Move DOWN by 'steps' (Library requires negative int)
		} else {
			migrationErr = m.Down() // Move DOWN completely
		}
	default:
		return errors.New("invalid migration mode: use 'up' or 'down'")
	}
	println(migrationErr)
	// Ignore ErrNoChange if nothing needed to be migrated
	if migrationErr != nil && !errors.Is(migrationErr, migrate.ErrNoChange) {
		fmt.Println("Database is already up to date!")
		return migrationErr
	}

	return nil
}

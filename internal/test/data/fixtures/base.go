package fixtures

import (
	"testing"

	"github.com/viettrung2103/bookmark-management/pkg/sqldb"
	"gorm.io/gorm"
)

type Fixture interface {
	//setup db
	SetupDB(db *gorm.DB)
	// migrate db
	Migrate() error
	//generate data
	GenerateData() error
	// get db
	DB() *gorm.DB
}

type base struct {
	db *gorm.DB
}

func (b *base) SetupDB(db *gorm.DB) {
	b.db = db
}

func (b *base) DB() *gorm.DB {
	return b.db
}

func NewFixture(t *testing.T, fix Fixture) *gorm.DB {
	//create test db
	fix.SetupDB(sqldb.CreateTestDb(t))

	// migrate schema
	err := fix.Migrate()
	if err != nil {
		t.Fatalf("Failed to migrate db: %v", err)
	}

	// generate data
	err = fix.GenerateData()
	if err != nil {
		t.Fatalf("Failed to generate data: %v", err)
	}

	// return db
	return fix.DB()
}

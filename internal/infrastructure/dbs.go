package infrastructure

import (
	"github.com/redis/go-redis/v9"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/common"
	redispkg "github.com/viettrung2103/bookmark-management/pkg/redis"
	"github.com/viettrung2103/bookmark-management/pkg/sqldb"
	"gorm.io/gorm"
)

// CreateRedisClient creates a new redis client
func CreateRedisClient() *redis.Client {
	redisClient, err := redispkg.NewClient("")
	common.HandleError(err)

	return redisClient
}

// CreateDBClient creates a new database client
func CreateDBClient() *gorm.DB {
	dbClient, err := sqldb.NewClient("")

	common.HandleError(err)
	err = MigrateDB(dbClient)
	common.HandleError(err)
	return dbClient
}

// MigrateDB migrates the database
func MigrateDB(sqlDB *gorm.DB) error {
	return sqlDB.AutoMigrate(&model.User{})

}

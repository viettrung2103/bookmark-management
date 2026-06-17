package main

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	api "github.com/viettrung2103/bookmark-management/internal/api"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/pkg/logger"
	redispkg "github.com/viettrung2103/bookmark-management/pkg/redis"
	"github.com/viettrung2103/bookmark-management/pkg/sqldb"
	"gorm.io/gorm"
)

// @title Bookmark API
// @version 2.0
// @description API for bookmark management
// @host localhost:8080
// @BasePath /
func main() {
	//create app config

	logger.SetLogLevel()

	// init app config
	cfg := createAPIConfig()

	// init redis
	redisClient := createRedisClient()

	//init db
	db := createDBClient()

	//init app
	app := createAPIApp(cfg, redisClient, db)

	//start app
	err := app.Start()
	if err != nil {
		panic(err)
	}

}

func createDBClient() *gorm.DB {
	dbClient, err := sqldb.NewClient("")
	if err != nil {
		panic(err)
	}
	dbClient.AutoMigrate(&model.User{})
	return dbClient
}

func createRedisClient() *redis.Client {
	redisClient, err := redispkg.NewClient("")
	if err != nil {
		panic(err)
	}
	return redisClient
}

func createAPIConfig() *config.Config {
	apiConfig, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	return apiConfig
}

func createAPIApp(cfg *config.Config, redis *redis.Client, db *gorm.DB) api.Engine {
	app := gin.New()
	a := api.NewEngine(app, cfg, redis, db)

	return a
}

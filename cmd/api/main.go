package main

import (
	api "github.com/viettrung2103/bookmark-management/internal/api"
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/pkg/logger"
	redispkg "github.com/viettrung2103/bookmark-management/pkg/redis"
)

// @title Bookmark API
// @version 1.5
// @description API for bookmark management
// @host localhost:8080
// @BasePath /
func main() {
	//create app config

	logger.SetLogLevel()

	//log.Debug().Str("name", "debug").Int("run-time", 1000).Msg("log nay chi hien thi o debug level")
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	redisClient, err := redispkg.NewClient("")
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg, redisClient)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}

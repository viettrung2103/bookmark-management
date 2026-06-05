package main

import (
	api "github.com/viettrung2103/bookmark-management/internal/api"
	"github.com/viettrung2103/bookmark-management/internal/config"
	redispkg "github.com/viettrung2103/bookmark-management/pkg/redis"
)

// @title Bookmark API
// @version 1.0
// @description API for bookmark management
// @host localhost:8080
// @BasePath /v1/links/
func main() {
	//create app config
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

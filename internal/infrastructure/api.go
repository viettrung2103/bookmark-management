package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/internal/api"
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/pkg/common"
)

func createAPIConfig() *config.Config {
	apiConfig, err := config.NewConfig()
	common.HandleError(err)

	return apiConfig
}

// func CreateAPIApp(cfg *config.Config, redis *redis.Client, db *gorm.DB) api.Engine {
func CreateAPIApp() api.Engine {

	app := gin.Default()
	//app := gin.New()
	//Engine: gin.New(),

	// init app config
	cfg := createAPIConfig()

	// init redis
	redisClient := CreateRedisClient()

	//init db
	db := CreateDBClient()

	//create jwt provider
	jwtGen, jwtVal := CreateJWTProvider()
	a := api.NewEngine(&api.EngineOpts{

		Engine: app,
		Cfg:    cfg,
		Redis:  redisClient,
		SqlDB:  db,
		JwtGen: jwtGen,
		JwtVal: jwtVal,
	})

	return a
}

func StartApp() {
	app := CreateAPIApp()
	err := app.Start()
	common.HandleError(err)

}

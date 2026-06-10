package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/viettrung2103/bookmark-management/docs"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"

	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/internal/handler"
	"github.com/viettrung2103/bookmark-management/internal/repository"
	"github.com/viettrung2103/bookmark-management/internal/service"
)

const version = 1

// Engine represents the application engine
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// engine struct implements Engine interface
type engine struct {
	app   *gin.Engine
	cfg   *config.Config
	redis *redis.Client
}

// NewEngine creates a new engine
func NewEngine(cfg *config.Config, redis *redis.Client) Engine {
	app := &engine{
		app:   gin.Default(),
		cfg:   cfg,
		redis: redis,
	}
	app.initRoutes()

	return app
}

// Start starts the engine
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTP implements the http.Handler interface to handle HTTP requests, to serve a specific request for testing purpose
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

// initRoutes initializes the routes
func (e *engine) initRoutes() {
	// genpass svc, handle and route

	shortenUrlRepo := repository.NewUrlStorage(e.redis)
	healthCheckRepo := repository.NewHealthCheck(e.redis)
	keyGen := stringutils.NewKeyGenerator()

	shortenUrlSvc := service.NewShortenUrl(shortenUrlRepo, keyGen)
	healthCheckSvc := service.NewHealthCheck(healthCheckRepo)

	shortenUrlHandler := handler.NewShortenLink(shortenUrlSvc, e.cfg)
	healthCheckHandler := handler.NewHealthCheck(healthCheckSvc)

	e.app.GET("/health-check", healthCheckHandler.CheckHealth)
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiPath := fmt.Sprintf("/v%d/links", version)

	apiBase := e.app.Group(apiPath)
	{
		// shorten link post
		apiBase.POST("/shorten", shortenUrlHandler.ShortenUrlLink)
		//apiBase.GET("/v1/links/shorten/:code", shortenUrlHandler.Redirect)

	}
}

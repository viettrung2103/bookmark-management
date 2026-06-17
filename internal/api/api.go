package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/viettrung2103/bookmark-management/docs"
	appHandler "github.com/viettrung2103/bookmark-management/internal/app/handler"
	userHandler "github.com/viettrung2103/bookmark-management/internal/app/handler/user"
	repo "github.com/viettrung2103/bookmark-management/internal/app/repository"
	userRepository "github.com/viettrung2103/bookmark-management/internal/app/repository/user"
	svc "github.com/viettrung2103/bookmark-management/internal/app/service"
	userService "github.com/viettrung2103/bookmark-management/internal/app/service/user"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"
	"gorm.io/gorm"

	"github.com/viettrung2103/bookmark-management/internal/config"
)

const version = 1

// Engine represents the application engine
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// engine struct implements Engine interface
type engine struct {
	eng   *gin.Engine
	cfg   *config.Config
	redis *redis.Client
	db    *gorm.DB
}

// NewEngine creates a new engine
func NewEngine(eng *gin.Engine, cfg *config.Config, redis *redis.Client, db *gorm.DB) Engine {
	app := &engine{
		//app:   gin.Default(),
		eng:   eng,
		cfg:   cfg,
		redis: redis,
		db:    db,
	}
	app.initRoutes()

	return app
}

// Start starts the engine
func (e *engine) Start() error {
	return e.eng.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTP implements the http.Handler interface to handle HTTP requests, to serve a specific request for testing purpose
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.eng.ServeHTTP(w, req)
}

type handlers struct {
	healthCheckHandler appHandler.HealthCheck
	linkHandler        appHandler.ShortenLink
	userHandler        userHandler.Handler
}

func (e *engine) initHandlers() *handlers {
	healthCheckRepo := repo.NewHealthCheck(e.redis)
	shortenUrlRepo := repo.NewUrlStorage(e.redis)

	keyGen := stringutils.NewKeyGenerator()

	shortenUrlSvc := svc.NewShortenUrl(shortenUrlRepo, keyGen)
	healthCheckSvc := svc.NewHealthCheck(healthCheckRepo)

	shortenUrlHdlr := appHandler.NewShortenLink(shortenUrlSvc, e.cfg)
	healthCheckHdlr := appHandler.NewHealthCheck(healthCheckSvc)

	userRepo := userRepository.NewRepository(e.db)
	userSvc := userService.NewService(userRepo)
	userHdlr := userHandler.NewHandler(userSvc)

	return &handlers{
		healthCheckHandler: healthCheckHdlr,
		linkHandler:        shortenUrlHdlr,
		userHandler:        userHdlr,
	}

}

// initRoutes initializes the routes
func (e *engine) initRoutes() {
	// genpass svc, handle and route

	//shortenUrlRepo := repo.NewUrlStorage(e.redis)
	//healthCheckRepo := repo.NewHealthCheck(e.redis)
	//keyGen := stringutils.NewKeyGenerator()
	//
	//shortenUrlSvc := svc.NewShortenUrl(shortenUrlRepo, keyGen)
	//healthCheckSvc := svc.NewHealthCheck(healthCheckRepo)
	//
	//shortenUrlHdlr := appHandler.NewShortenLink(shortenUrlSvc, e.cfg)
	//healthCheckHandler := appHandler.NewHealthCheck(healthCheckSvc)

	allHandlers := e.initHandlers()

	e.eng.GET("/health-check", allHandlers.healthCheckHandler.CheckHealth)

	//int swagger routes
	docs.SwaggerInfo.Host = e.cfg.Hostname
	e.eng.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiPath := fmt.Sprintf("/v%d", version)

	apiBase := e.eng.Group(apiPath)
	{
		// link route
		linkBase := apiBase.Group("/links")

		linkBase.POST("/shorten", allHandlers.linkHandler.ShortenUrlLink)
		linkBase.GET("/redirect/:code", allHandlers.linkHandler.RedirectUrl)

		//user route
		userBase := apiBase.Group("/users")
		userBase.POST("/register", allHandlers.userHandler.Register)

	}
}

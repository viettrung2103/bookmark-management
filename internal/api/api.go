package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/viettrung2103/bookmark-management/docs"
	"github.com/viettrung2103/bookmark-management/internal/api/middleware"
	bookmarkHandler "github.com/viettrung2103/bookmark-management/internal/app/handler/bookmark"
	healthCheckHandler "github.com/viettrung2103/bookmark-management/internal/app/handler/healthcheck"
	urlHandler "github.com/viettrung2103/bookmark-management/internal/app/handler/url"
	userHandler "github.com/viettrung2103/bookmark-management/internal/app/handler/user"

	bookmarkRepository "github.com/viettrung2103/bookmark-management/internal/app/repository/bookmark"
	healthCheckRepository "github.com/viettrung2103/bookmark-management/internal/app/repository/healthcheck"
	urlRepository "github.com/viettrung2103/bookmark-management/internal/app/repository/url"
	userRepository "github.com/viettrung2103/bookmark-management/internal/app/repository/user"

	bookmarkService "github.com/viettrung2103/bookmark-management/internal/app/service/bookmark"
	healthCheckService "github.com/viettrung2103/bookmark-management/internal/app/service/healthcheck"
	urlService "github.com/viettrung2103/bookmark-management/internal/app/service/url"
	userService "github.com/viettrung2103/bookmark-management/internal/app/service/user"

	"github.com/viettrung2103/bookmark-management/pkg/jwtutils"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"
	"gorm.io/gorm"

	"github.com/viettrung2103/bookmark-management/internal/config"
)

const version = 1

// Engine represents the application engine
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	InitRoutes()
}

// engine struct implements Engine interface
type engine struct {
	eng    *gin.Engine
	cfg    *config.Config
	redis  *redis.Client
	db     *gorm.DB
	jwtGen jwtutils.JWTGenerator
	jwtVal jwtutils.JWTValidator
}

// EngineOpts holds initialization dependencies for the engine
type EngineOpts struct {
	Engine *gin.Engine
	Cfg    *config.Config
	Redis  *redis.Client
	SqlDB  *gorm.DB
	JwtGen jwtutils.JWTGenerator
	JwtVal jwtutils.JWTValidator
}

// New creates a new engine
func NewEngine(opts *EngineOpts) Engine {
	app := &engine{
		eng:    opts.Engine,
		cfg:    opts.Cfg,
		redis:  opts.Redis,
		db:     opts.SqlDB,
		jwtGen: opts.JwtGen,
		jwtVal: opts.JwtVal,
	}
	app.InitRoutes()

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
	healthCheckHandler healthCheckHandler.Handler
	linkHandler        urlHandler.Handler
	userHandler        userHandler.Handler
	bookmarkHandler    bookmarkHandler.Handler
}

func (e *engine) initHandlers() *handlers {
	healthCheckRepo := healthCheckRepository.NewRepository(e.redis)
	shortenUrlRepo := urlRepository.NewRepository(e.redis)

	keyGen := stringutils.NewKeyGenerator()

	shortenUrlSvc := urlService.NewService(shortenUrlRepo, keyGen)
	healthCheckSvc := healthCheckService.NewService(healthCheckRepo)

	shortenUrlHdlr := urlHandler.NewShortenLink(shortenUrlSvc, e.cfg)
	healthCheckHdlr := healthCheckHandler.NewHandler(healthCheckSvc)

	passwordHashing := stringutils.NewPasswordHasher()

	userRepo := userRepository.NewRepository(e.db)
	//userSvc := userService.NewService(userRepo)
	//userHdlr := userHandler.NewHandler(userSvc)

	userSvcInput := &userService.UserServiceOpts{
		UserRepo:        userRepo,
		PasswordHashing: passwordHashing,
		JwtGenerator:    e.jwtGen,
	}

	userSvc := userService.NewService(userSvcInput)
	userHdlr := userHandler.NewHandler(userSvc)

	bookmarkRepo := bookmarkRepository.NewRepository(e.db)
	bookmarkSvcOpts := &bookmarkService.BookmarkServiceOpts{
		Keygen:             keyGen,
		BookmarkRepository: bookmarkRepo,
	}

	bookmarkSvc := bookmarkService.NewService(bookmarkSvcOpts)
	bookmarkHdlr := bookmarkHandler.NewHandler(bookmarkSvc)

	return &handlers{
		healthCheckHandler: healthCheckHdlr,
		linkHandler:        shortenUrlHdlr,
		userHandler:        userHdlr,
		bookmarkHandler:    bookmarkHdlr,
	}

}

// initRoutes initializes the routes
func (e *engine) InitRoutes() {

	allHandlers := e.initHandlers()
	jwtAuth := middleware.NewJWTAuth(e.jwtVal)

	e.eng.GET("/health-check", allHandlers.healthCheckHandler.HealthCheck)

	//int swagger routes
	docs.SwaggerInfo.Host = e.cfg.Hostname
	e.eng.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiRoute := fmt.Sprintf("/v%d", version)

	OpenBase := e.eng.Group(apiRoute)
	{
		// url route
		linkBase := OpenBase.Group("/links")

		linkBase.POST("/shorten", allHandlers.linkHandler.ShortenUrlLink)
		linkBase.GET("/redirect/:code", allHandlers.linkHandler.RedirectUrl)

		//user route
		userBase := OpenBase.Group("/users")
		userBase.POST("/register", allHandlers.userHandler.Register)
		userBase.POST("/login", allHandlers.userHandler.Login)

	}

	privateBase := e.eng.Group(apiRoute) // /v1
	privateBase.Use(jwtAuth.JWTAuthMiddleWare())
	{
		// user-related
		privateBase.GET("/self/info", allHandlers.userHandler.SelfInfo)
		privateBase.PUT("/self/info", allHandlers.userHandler.EditSelfInfo)

		//bookmark
		privateBase.POST("/bookmarks", allHandlers.bookmarkHandler.AddBookmark)
	}

	//test
	OpenBase.GET("/test/:id/blah/:uid", userHandler.Test)
}

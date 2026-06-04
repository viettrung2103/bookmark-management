package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/internal/handler"
	"github.com/viettrung2103/bookmark-management/internal/service"
)

// Engine represents the application engine
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// engine struct implements Engine interface
type engine struct {
	app *gin.Engine
	cfg *config.Config
}

// NewEngine creates a new engine
func NewEngine(cfg *config.Config) Engine {
	app := &engine{
		app: gin.Default(),
		cfg: cfg,
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
	genIdSvc := service.NewGenId()
	genIdHanlder := handler.NewGenId(genIdSvc, e.cfg)

	e.app.GET("/health-check", genIdHanlder.GenerateId)

}

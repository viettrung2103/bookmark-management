package url

import (
	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/internal/app/service/urlstorage"
	"github.com/viettrung2103/bookmark-management/internal/config"
)

// ShortenLink represents the shorten url handler
type Handler interface {
	ShortenUrlLink(c *gin.Context)
	//CheckHealth(c *gin.Context)
	RedirectUrl(c *gin.Context)
}
type shortenLinkHandler struct {
	shortenLinkService urlstorage.Service
	cfg                *config.Config
}

// NewShortenLink creates a new ShortenLink
func NewShortenLink(shortenLinkSvc urlstorage.Service, cfg *config.Config) Handler {
	return &shortenLinkHandler{
		shortenLinkService: shortenLinkSvc,
		cfg:                cfg,
	}
}

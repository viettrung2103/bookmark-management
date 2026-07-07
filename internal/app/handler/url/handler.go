package url

import (
	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/internal/app/service/url"
	"github.com/viettrung2103/bookmark-management/internal/config"
)

// ShortenLink represents the shorten url handler
type Handler interface {
	ShortenUrlLink(c *gin.Context)
	//CheckHealth(c *gin.Context)
	RedirectUrl(c *gin.Context)
}
type shortenLinkHandler struct {
	shortenLinkService link.URLService
	cfg                *config.Config
}

// NewShortenLink creates a new ShortenLink
func NewShortenLink(shortenLinkSvc link.URLService, cfg *config.Config) Handler {
	return &shortenLinkHandler{
		shortenLinkService: shortenLinkSvc,
		cfg:                cfg,
	}
}

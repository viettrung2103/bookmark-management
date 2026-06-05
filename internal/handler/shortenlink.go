package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/internal/config"
	"github.com/viettrung2103/bookmark-management/internal/service"
)

// ShortenLink represents the shorten url handler
type ShortenLink interface {
	ShortenUrlLink(c *gin.Context)
	//CheckHealth(c *gin.Context)
	RedirectUrl(c *gin.Context)
}
type shortenLinkHandler struct {
	shortenLinkService service.ShortenUrl
	cfg                *config.Config
}

// ShortenRequest represents the shorten url request
type ShortenRequest struct {
	Url              string `json:"url" binding:"required,url"`
	ExpiringDuration int    `json:"exp" binding:"required"`
}

// NewShortenLink creates a new ShortenLink
func NewShortenLink(shortenLinkSvc service.ShortenUrl, cfg *config.Config) ShortenLink {
	return &shortenLinkHandler{
		shortenLinkService: shortenLinkSvc,
		cfg:                cfg,
	}
}

// ShortenUrlLink shortens a url
// @Summary shorten a linkn
// @Description receive a link, return a code
// @Tags shorten-url
// @Success 200 {object} map[string]interface{}
// @Router /v1/links/shorten [post]
func (h *shortenLinkHandler) ShortenUrlLink(c *gin.Context) {
	var req ShortenRequest

	// bind the incoming request with our struct
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	code, err := h.shortenLinkService.ShortenUrlWithExpiringTime(c, req.Url, req.ExpiringDuration)
	if err != nil {

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Err"})
		return

	}

	c.JSON(http.StatusOK,
		gin.H{
			"code":    code,
			"message": "Shorten URL generated successfully",
		})
}

func (h *shortenLinkHandler) RedirectUrl(c *gin.Context) {

}

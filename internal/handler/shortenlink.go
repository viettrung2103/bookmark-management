package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

// shortenUrlRequest represents the shorten url request
type shortenUrlRequest struct {
	Url              string `json:"url" binding:"required,url"`
	ExpiringDuration int    `json:"exp" binding:"required"`
}

type shortenUrlResponse struct {
	Code string `json:"code"`
}

// NewShortenLink creates a new ShortenLink
func NewShortenLink(shortenLinkSvc service.ShortenUrl, cfg *config.Config) ShortenLink {
	return &shortenLinkHandler{
		shortenLinkService: shortenLinkSvc,
		cfg:                cfg,
	}
}

// ShortenUrlLink shorten the url to code
// @Summary receive the url, return the code
// @Tags link
// @Accept application/json
// @Produce application/json
// @Param request body shortenUrlRequest true "Shorten URL Input payload"
// @Success 200 {object} string
// @Router /v1/links/shorten [post]
func (h *shortenLinkHandler) ShortenUrlLink(c *gin.Context) {
	var req shortenUrlRequest

	// bind the incoming request with our struct
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.shortenurl.ShortenUrlLink").Msg("failed to get req body from code")

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

// Redirect Forward the request to the original url
// @Summary Redirect Forward the request to the original url
// @Tags link
// @Accept application/json
// @Produce application/json
// @Param code path string true "code"
// @Success	302
// @Router /v1/links/shorten/{code} [get]
func (h *shortenLinkHandler) RedirectUrl(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "code is required"})
	}

	// call service to get url from code
	url, err := h.shortenLinkService.GetLinkFromCode(c, code)
	if err != nil {
		if errors.Is(err, service.ErrCodeDoesNotExist) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url not found"})
			return
		}
		log.Error().Err(err).Str("from", "handler.shortenurl.Redirect").Msg("failed to get url from code")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// redirect to url
	c.Redirect(http.StatusFound, url)

}

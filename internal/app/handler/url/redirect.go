package url

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/viettrung2103/bookmark-management/internal/app/service/url"
)

// Redirect Forward the request to the original url
// @Summary Redirect Forward the request to the original url
// @Tags url
// @Accept application/json
// @Produce application/json
// @Param code path string true "code"
// @Success	302
// @Router /v1/links/redirect/{code} [get]
func (h *shortenLinkHandler) RedirectUrl(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "code is required"})
	}

	// call service to get url from code
	url, err := h.shortenLinkService.GetLinkFromCode(c, code)
	if err != nil {
		if errors.Is(err, link.ErrCodeDoesNotExist) {
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

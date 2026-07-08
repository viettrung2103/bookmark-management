package url

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/pkg/requestutils"
)

// shortenUrlRequest represents the shorten url request
type shortenUrlRequest struct {
	URL              string `json:"url" binding:"required,url"`
	ExpiringDuration int    `json:"exp" binding:"required"`
}

type shortenUrlResponse struct {
	Code string `json:"code"`
}

// ShortenUrlLink shorten the url to code
// @Summary receive the url, return the code
// @Tags url
// @Accept application/json
// @Produce application/json
// @Param request body shortenUrlRequest true "Shorten URL Input payload"
// @Success 200 {object} string
// @Router /v1/links/shorten [post]
func (h *shortenLinkHandler) ShortenUrlLink(c *gin.Context) {
	//var req shortenUrlRequest

	// bind the incoming request with our struct
	//err := c.ShouldBindJSON(&req)
	//if err != nil {
	//	log.Error().Err(err).Str("from", "handler.shortenurl.ShortenUrlLink").Msg("failed to get req body from code")
	//
	//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	//	return
	//}
	input, err := requestutils.BindInputFromRequest[shortenUrlRequest](c)
	if err != nil {
		return
	}

	code, err := h.shortenLinkService.ShortenUrlWithExpiringTime(c, input.URL, input.ExpiringDuration)
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

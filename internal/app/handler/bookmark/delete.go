package bookmark

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/viettrung2103/bookmark-management/pkg/common"
	"github.com/viettrung2103/bookmark-management/pkg/requestutils"
	"github.com/viettrung2103/bookmark-management/pkg/response"
)

// DeleteBookmark delete bookmark of id of current user
// @Summary delete bookmark of id of current user
// @Description delete bookmark of id of current user
// @Tags bookmark
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} object{data=[]model.Bookmark}
// @Router /v1/bookmarks/{id} [delete]
func (h *bookmarkHandler) DeleteBookmark(c *gin.Context) {
	id := c.Param("id")

	userID, err := requestutils.GetUserIDFromRequest(c)
	common.HandleError(err)
	err = h.svc.DeleteBookmarkByID(c, userID, id)

	switch {
	case errors.Is(err, nil):
		break
	default:
		log.Error().Err(err).Str("userID", userID).Msg("GetBookmarks err")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})

}

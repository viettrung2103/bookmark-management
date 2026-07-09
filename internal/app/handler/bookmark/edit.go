package bookmark

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	"github.com/viettrung2103/bookmark-management/pkg/requestutils"
	"github.com/viettrung2103/bookmark-management/pkg/response"
)

//type loginInput struct {
//	Username string `json:"username" binding:"required"`
//	Password string `json:"password" binding:"required,gte=8"`
//}

type putBookmarkInput struct {
	Description string `json:"description" binding:"required"`
	URL         string `json:"url" binding:"required,url"`
	//ID          int ``
}

// @Summary edit bookmark of current user
// @Description edit bookmark of current user
// @Tags bookmark
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Param input body putBookmarkInput true "Input required"
// @Param id path string true "id"
// @Success 200 {object} object{data=[]model.Bookmark}
// @Router /v1/bookmarks/{id} [put]
func (h *bookmarkHandler) EditBookmark(c *gin.Context) {
	// get bookmark id from request endpoint
	id := c.Param("id")
	// get input from request : input description and url, user id
	input, uid, err := requestutils.BindInputFromRequestWithAuth[putBookmarkInput](c)
	if err != nil {
		return
	}

	err = h.svc.EditBookmarkByID(c, uid, id, input.Description, input.URL)

	switch {
	case errors.Is(err, dbutils.ErrDuplication):
		c.JSON(http.StatusBadRequest, &response.Message{
			Message: err.Error(),
		})
		return
	case errors.Is(err, nil):
		break
	default:
		log.Err(err).Str("bookmarkID", id).Msg("failed to update bookmark info")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

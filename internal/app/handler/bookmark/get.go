package bookmark

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/viettrung2103/bookmark-management/internal/app/handler/dto"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/requestutils"
	"github.com/viettrung2103/bookmark-management/pkg/response"
)

type getBookmarkInput struct {
	Page  int `form:"page" validate:"gte=1"`
	Limit int `form:"limit" validate:"gte=1"`
}

// @Summary get lists of bookmarks of current user
// @Description get lists of bookmarks of current user
// @Tags bookmark
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Param page query integer false "page" Format(int32) default(1)
// @Param limit query integer false "limit" Format(int32) default(20)
// @Success 200 {object} object{data=[]model.Bookmark}
// @Router /v1/bookmarks [get]
func (h *bookmarkHandler) GetBookmarks(c *gin.Context) {
	input, uid, err := requestutils.BindInputFromRequestWithAuth[getBookmarkInput](c)
	if err != nil {
		return
	}

	res, err := h.svc.GetBookmarks(c, uid, input.Page, input.Limit)
	switch {
	case errors.Is(err, nil):
		break
	default:
		log.Error().Err(err).Str("userID", uid).Msg("GetBookmarks err")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, &dto.SuccessResponse[[]*model.Bookmark]{
		Data: res.Bookmarks,
		Pagination: &dto.Pagination{
			Page:  input.Page,
			Limit: input.Limit,
			Total: res.Count,
		},
	})
}

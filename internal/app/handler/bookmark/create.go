package bookmark

import (
	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/pkg/requestutils"
)

type addBookmarkInput struct {
	Description string `json:"description" example:"Google" validate:"lte=255"`
	URL         string `json:"url" example:"https://www.google.com" validate:"required,url,lte=2048"`
}

// AddBookmark Create a new bookmark for a bookmark group
// @Summary Create a new bookmark for a bookmark group
// @Description Create a new bookmark for a bookmark group
// @Tags bookmark
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Param input body addBookmarkInput true "Input required"
// @Success 200 {object} object{data=model.Bookmark,message=string} "Success"
// @Router /v1/bookmarks [post]
func (h *bookmarkHandler) AddBookmark(c *gin.Context) {
	// lay input
	println("Getting input")
	input, uid, err := requestutils.BindInputFromRequestWithAuth[addBookmarkInput](c)
	if err != nil {
		return
	}
	println("add bookmark using service")

	// call service
	newBookmark, err := h.svc.AddBookmark(c, input.Description, input.URL, uid)
	if err != nil {
		return
	}
	c.JSON(200, gin.H{
		"data":    newBookmark,
		"message": "Create a bookmark successfully",
	})

	// tra ve response
}

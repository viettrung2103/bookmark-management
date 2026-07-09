package bookmark

import (
	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/internal/app/service/bookmark"
)

// Handler interface for user handler
type Handler interface {
	AddBookmark(c *gin.Context)
	GetBookmarks(c *gin.Context)
}

type bookmarkHandler struct {
	svc bookmark.Service
}

// NewHandler creates a new user handler
func NewHandler(bookmarkSvc bookmark.Service) Handler {
	return &bookmarkHandler{svc: bookmarkSvc}
}

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/internal/app/service/user"
)

// Handler interface for user handler
type Handler interface {
	Register(c *gin.Context)
}

type userHandler struct {
	service user.Service
}

// NewHandler creates a new user handler
func NewHandler(service user.Service) Handler {
	return &userHandler{service: service}
}

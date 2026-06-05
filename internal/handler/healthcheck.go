package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/internal/service"
)

type HealthCheck interface {
	CheckHealth(c *gin.Context)
}
type healthCheckHandler struct {
	healthCheckSvc service.HealthCheck
	//cfg                *config.Config
}

func NewHealthCheck(healthCheckSvc service.HealthCheck) HealthCheck {
	return &healthCheckHandler{
		healthCheckSvc: healthCheckSvc,
	}
}

// CheckHealth checks the health of the service
// @Summary check redis health
// @Description ping and pong with redis server
// @Tags shorten-url
// @Success 200 {object} map[string]interface{}
// @Router /check-health [get]
func (h *healthCheckHandler) CheckHealth(c *gin.Context) {
	err := h.healthCheckSvc.CheckHealth(c)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
			"redis":  "unreachable",
			"error":  "Service Unavailable",
		})
		return
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status": "UP",
			"redis":  "reachable",
		})
}

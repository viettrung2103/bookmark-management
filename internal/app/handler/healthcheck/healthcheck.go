package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// CheckHealth checks the health of the service
// @Summary check redis health
// @Description ping and pong with redis server
// @Tags health check
// @Success 200 {object} map[string]interface{}
// @Router /health-check [get]
func (h *healthCheckHandler) CheckHealth(c *gin.Context) {
	err := h.healthCheckSvc.HealthCheck(c)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.healthCheckHandler.CheckHealth").Msg("redis server is down")

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

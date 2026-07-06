package healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/internal/app/service/healthcheck"
)

// HealthCheck interface for health check
type Handler interface {
	HealthCheck(c *gin.Context)
}
type healthCheckHandler struct {
	healthCheckSvc healthcheck.HealthCheckService
}

// NewHealthCheck creates a new health check handler
func NewHandler(healthCheckSvc healthcheck.HealthCheckService) Handler {
	return &healthCheckHandler{
		healthCheckSvc: healthCheckSvc,
	}
}

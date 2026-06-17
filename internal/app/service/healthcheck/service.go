package healthcheck

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/repository/healthcheck"
)

// HealthCheck represents the health check service
//
//go:generate mockery --name=HealthCheck --filename=healthcheck.go
type Service interface {
	HealthCheck(ctx context.Context) error
}

type healthCheckService struct {
	repo healthcheck.Repository
}

// NewHealthCheckS creates a new HealthCheck
func NewService(repo healthcheck.Repository) Service {
	return &healthCheckService{repo: repo}
}

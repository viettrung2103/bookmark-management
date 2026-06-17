package service

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/repository"
)

// HealthCheck represents the health check service
//
//go:generate mockery --name=HealthCheck --filename=healthcheck.go
type HealthCheck interface {
	CheckHealth(ctx context.Context) error
}

type healthCheckService struct {
	repo repository.HealthCheck
}

// NewHealthCheckS creates a new HealthCheck
func NewHealthCheck(repo repository.HealthCheck) HealthCheck {
	return &healthCheckService{repo: repo}
}

// CheckHealth checks the health of the service
func (s *healthCheckService) CheckHealth(ctx context.Context) error {
	return s.repo.CheckHealth(ctx)
}

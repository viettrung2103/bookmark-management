package healthcheck

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// HealthCheck is the interface for health check
//
//go:generate mockery --name=HealthCheckRepository --filename=../../mocks/healthcheck.go
type HealthCheckRepository interface {
	HealthCheck(ctx context.Context) error
}
type healthCheckRepo struct {
	c *redis.Client
}

// NewHealthCheck creates a new HealthCheck
func NewRepository(c *redis.Client) HealthCheckRepository {
	return &healthCheckRepo{
		c: c,
	}
}

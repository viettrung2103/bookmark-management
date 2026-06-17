package healthcheck

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// HealthCheck is the interface for health check
//
//go:generate mockery --name=HealthCheck --filename=healthcheck.go
type Repository interface {
	HealthCheck(ctx context.Context) error
}
type healthCheckRepo struct {
	c *redis.Client
}

// NewHealthCheck creates a new HealthCheck
func NewRepository(c *redis.Client) Repository {
	return &healthCheckRepo{
		c: c,
	}
}

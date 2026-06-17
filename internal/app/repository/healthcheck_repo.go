package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// HealthCheck is the interface for health check
//go:generate mockery --name=HealthCheck --filename=healthcheck.go

type HealthCheck interface {
	CheckHealth(ctx context.Context) error
}
type healthCheckRepo struct {
	c *redis.Client
}

// NewHealthCheck creates a new HealthCheck
func NewHealthCheck(c *redis.Client) HealthCheck {
	return &healthCheckRepo{
		c: c,
	}
}

// CheckHealth check health of redis server
func (s *healthCheckRepo) CheckHealth(ctx context.Context) error {

	return s.c.Ping(ctx).Err()
}

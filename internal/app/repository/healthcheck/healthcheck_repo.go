package healthcheck

import (
	"context"
)

// CheckHealth check health of redis server
func (s *healthCheckRepo) HealthCheck(ctx context.Context) error {

	return s.c.Ping(ctx).Err()
}

package healthcheck

import (
	"context"
)

// CheckHealth checks the health of the service
func (s *healthCheckService) HealthCheck(ctx context.Context) error {
	return s.repo.HealthCheck(ctx)
}

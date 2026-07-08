package healthcheck

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/viettrung2103/bookmark-management/pkg/redis"
)

func TestHealthCheckRepo_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupClient func(client *redis.Client)
		wantErr     bool // Simply check if an error occurred
	}{
		{
			name: "success - redis is up",
			setupClient: func(client *redis.Client) {
				// The miniredis server is running and healthy. Do nothing.
			},
			wantErr: false,
		},
		{
			name: "failure - redis connection broken",
			setupClient: func(client *redis.Client) {
				// Forcefully close the client to simulate a connection outage
				client.Close()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// 1. Initialize the in-memory Redis server using your helper
			client := redisPkg.InitMockRedis(t)

			// 2. Manipulate the client state based on the test case
			tc.setupClient(client)

			// 3. Initialize your repository
			repo := NewRepository(client)

			// 4. Execute the HealthCheck
			err := repo.HealthCheck(context.Background())

			// 5. Assert the results
			if tc.wantErr {
				assert.Error(t, err) // We expect some kind of network/closed client error
			} else {
				assert.NoError(t, err) // We expect a successful PONG (nil error)
			}
		})
	}
}

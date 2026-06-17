package jwtutils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewJWTGenerator tests the NewJWTGenerator function
func TestNewJWTGenerator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		keyPath string

		expectedErrStr string
	}{
		{
			name:           "normal case",
			keyPath:        filepath.FromSlash("./test.private.pem"),
			expectedErrStr: "",
		},
		{
			name:           "err case - file not found",
			keyPath:        filepath.FromSlash("./non-exist.pem"),
			expectedErrStr: "open",
		},
		{
			name:           "err case - not a private key",
			keyPath:        filepath.FromSlash("./test.public.pem"),
			expectedErrStr: "structure error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewJWTGenerator(tc.keyPath)
			if err != nil {
				assert.ErrorContains(t, err, tc.expectedErrStr)
			}
		})
	}
}

// TestGenerateJWT tests the GenerateJWT method
func TestGenerateJWT(t *testing.T) {
	expectedToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0In0.OSrRquH2kdcpUWZu9YyWPj-iUZIYBNNi_mkZ3X6D1D9UYcFKJZFBUS-h0KEHOyhs3yGpEXqYBa1sC7cw3NDdmD1KVUh-LunkNX3swo8Dydi8HcXth6V6k8ztqUle-X5lh9MVsaVT9nMliW-UbCwceJugZmqdZoSEKYwmv9gyQX8WouHonglAtnrE7ejz8pIf04Gq0lSTQuApgmPgwhKbH7fUAzW2QLD_O7HGBsxzauaz9S0UWSNck2xGVZxPSMnw2CtsocJ6m0ZEUwnJivomwt_LnNLjD2B5k3dHU8b_FQ0aqyaIfs9fNfDZZCeosamMpKvYmDYWa24mJoUJ6jwYqg"
	t.Parallel()
	gen, err := NewJWTGenerator(filepath.FromSlash("./test.private.pem"))
	if err != nil {
		t.Fatal("should not fail")
	}

	token, err := gen.GenerateJWT(map[string]any{
		"sub": "1234",
	})

	if err != nil {
		t.Fatal("should not fail")
	}
	assert.Equal(t, expectedToken, token)

}

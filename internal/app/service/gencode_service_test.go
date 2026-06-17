package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		expectedLength int
	}{
		{
			name:           "success",
			expectedLength: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewCode()
			password := testSvc.GenerateCode()
			assert.Equal(t, tc.expectedLength, len(password))
		})
	}
}

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testUUID = "123e4567-e89b-12d3-a456-426614174000"

func TestGeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		expectedLength int
	}{
		{
			name:           "success",
			expectedLength: len(testUUID),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewGenId()
			password := testSvc.GenerateId()
			//assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedLength, len(password))
		})
	}
}

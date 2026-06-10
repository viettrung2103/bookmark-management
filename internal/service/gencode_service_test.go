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
		//expectedError  error
	}{
		{
			name:           "success",
			expectedLength: 10,
			//expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewCode()
			password := testSvc.GenerateCode()
			//assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedLength, len(password))
		})
	}
}

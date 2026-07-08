package stringutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHasher_Hashing(t *testing.T) {
	t.Parallel()

	hasher := NewPasswordHasher()

	testCases := []struct {
		name        string
		inputString string
	}{
		{
			name:        "success - normal password",
			inputString: "supersecret123",
		},
		{
			name:        "success - short password",
			inputString: "123",
		},
		{
			name:        "success - empty password",
			inputString: "", // bcrypt handles empty strings without panicking
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// 1. Generate the hash
			hashed := hasher.Hashing(tc.inputString)

			// 2. Assert it is not empty and not plain text
			assert.NotEmpty(t, hashed)
			assert.NotEqual(t, tc.inputString, hashed, "Hashed password should not equal the plain text")

			// 3. Verify it is a valid bcrypt hash by comparing it back to the input
			match := hasher.CompareHashedPassword(hashed, tc.inputString)
			assert.True(t, match, "The hashed password should match the original input")
		})
	}
}

func TestPasswordHasher_CompareHashedPassword(t *testing.T) {
	t.Parallel()

	hasher := NewPasswordHasher()

	// Create a real hash to test against in our scenarios
	validPassword := "mysecurepassword"
	validHash := hasher.Hashing(validPassword)

	testCases := []struct {
		name          string
		inputHash     string
		inputPassword string
		expectedMatch bool
	}{
		{
			name:          "success - passwords match",
			inputHash:     validHash,
			inputPassword: validPassword, // Matches the one used to generate validHash
			expectedMatch: true,
		},
		{
			name:          "failure - wrong password",
			inputHash:     validHash,
			inputPassword: "wrongpassword",
			expectedMatch: false,
		},
		{
			name:          "failure - invalid hash format",
			inputHash:     "this-is-not-a-valid-bcrypt-string",
			inputPassword: validPassword,
			expectedMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			match := hasher.CompareHashedPassword(tc.inputHash, tc.inputPassword)

			assert.Equal(t, tc.expectedMatch, match)
		})
	}
}

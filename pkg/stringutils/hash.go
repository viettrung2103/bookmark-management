package stringutils

import "golang.org/x/crypto/bcrypt"

//go:generate mockery --name=PasswordHashing --filename=hashing.go
// PasswordHashing interface for password hashing

type PasswordHashing interface {
	Hashing(input string) string
	CompareHashedPassword(hashedString, input string) bool
}

type passwordHasher struct{}

// NewPasswordHasher creates a new password hasher
func NewPasswordHasher() PasswordHashing {
	return &passwordHasher{}
}

// Hashing hashes the input string
func (h *passwordHasher) Hashing(input string) string {
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	return string(hashedBytes)
}

// CompareHashedPassword compares the hashed password with the input password
func (h *passwordHasher) CompareHashedPassword(hashedString, input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(input))
	return err == nil
}

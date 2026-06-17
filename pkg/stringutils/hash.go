package stringutils

import "golang.org/x/crypto/bcrypt"

// Hastring hashes the input string
func Hastring(input string) string {

	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	return string(hashedBytes)
}

// CompareHashedPassword compares the hashed password with the input password
func CompareHashedPassword(hashedString, input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(input))
	return err == nil
}

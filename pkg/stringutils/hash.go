package stringutils

import "golang.org/x/crypto/bcrypt"

func Hastring(input string) string {

	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	return string(hashedBytes)
}

func CompareHashedPassword(hashedString, input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(input))
	return err == nil
}

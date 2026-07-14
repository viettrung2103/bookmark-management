package jwtutils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtExpireDuration = 1 * time.Hour
)

var ErrInvalidCreditials = errors.New("invalid username or password")

// TokenInfo contains user ID and username
type TokenInfo struct {
	Id       string
	Username string
}

func GetMapClaim(uid, username string) jwt.MapClaims {
	return jwt.MapClaims{
		"uid":      uid,
		"username": username,
		"exp":      time.Now().Add(jwtExpireDuration).Unix(),
		"iat":      time.Now().Unix(),
	}
}

func (t *TokenInfo) ToMapClaim() jwt.MapClaims {
	return GetMapClaim(t.Id, t.Username)
}

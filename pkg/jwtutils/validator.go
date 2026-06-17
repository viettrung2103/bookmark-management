package jwtutils

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExtractToken = errors.New("failed to extract token")
)

type JWTValidator interface {
	ValidateJWT(tokenStr string) (jwt.MapClaims, error)
}

type validator struct {
	publicKey *rsa.PublicKey
}

func NewJWTValidator(publicKeyPath string) (JWTValidator, error) {
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	return &validator{
		publicKey: publicKey,
	}, nil

}

func (v *validator) ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return v.publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, ErrExtractToken
}

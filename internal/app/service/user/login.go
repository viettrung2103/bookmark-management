package user

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
)

const jwtExpireDuration = 1 * time.Hour

var ErrInvalidCreditials = errors.New("invalid username or password")

// TokenInfo contains user ID and username
type TokenInfo struct {
	Id       string
	Username string
}

// ToMapClaim converts TokenInfo to jwt.MapClaims
func (t *TokenInfo) ToMapClaim() jwt.MapClaims {
	return jwt.MapClaims{
		"uid":      t.Id,
		"username": t.Username,
		"exp":      time.Now().Add(jwtExpireDuration).Unix(),
		"iat":      time.Now().Unix(),
	}
}

// Login logs in a user and returns a JWT
func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
	// get user tu username
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, dbutils.ErrRecordNotFound) {
			return "", ErrInvalidCreditials
		}
		return "", err
	}

	// check password --> compare hash with passwor
	check := s.passwordHashing.CompareHashedPassword(user.Password, password)
	if !check {
		return "", ErrInvalidCreditials
	}

	// generate token
	tokenInfo := &TokenInfo{
		Id:       user.ID,
		Username: user.Username,
	}
	token, err := s.jwtGenerator.GenerateJWT(tokenInfo.ToMapClaim())
	if err != nil {
		return "", err
	}
	// return token
	return token, nil
}

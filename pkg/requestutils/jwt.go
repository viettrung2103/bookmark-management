package requestutils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/viettrung2103/bookmark-management/pkg/response"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrInvalidUID   = errors.New("invalid user")
)

// GetJWTClaimsFromRequest gets JWT claims from request
func GetJWTClaimsFromRequest(c *gin.Context) (jwt.MapClaims, error) {
	tokenInfo, _ := c.Get("claims")
	claims, valid := tokenInfo.(jwt.MapClaims)
	if !valid {
		c.JSON(http.StatusUnauthorized, &response.Message{
			Message: "invalid jwt token",
		})
		c.Abort()
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// GetUserIDFromRequest gets user ID from request
func GetUserIDFromRequest(c *gin.Context) (string, error) {
	claims, err := GetJWTClaimsFromRequest(c)
	if err != nil {
		return "", err
	}

	uid, ok := claims["uid"].(string)

	if !ok || uid == "" {
		c.JSON(http.StatusUnauthorized, &response.Message{
			Message: "invalid jwt token",
		})
		c.Abort()
		return "", ErrInvalidUID
	}
	return uid, nil
}

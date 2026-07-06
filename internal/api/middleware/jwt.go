package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/viettrung2103/bookmark-management/pkg/jwtutils"
	"github.com/viettrung2103/bookmark-management/pkg/response"
)

type JWTAuth interface {
	JWTAuthMiddleWare() gin.HandlerFunc
}

type jwtAuth struct {
	jwtValidator jwtutils.JWTValidator
}

func NewJWTAuth(jwtValidator jwtutils.JWTValidator) JWTAuth {
	return &jwtAuth{
		jwtValidator: jwtValidator,
	}
}

func (j *jwtAuth) JWTAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get auth header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Message{Message: "Unauthorized"})
			c.Abort()
			return
		}

		// extract token tu header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.Message{Message: "Auth header has invalid format"})
			c.Abort()
			return
		}
		// get tokenClaims
		token := parts[1]

		tokenClaims, err := j.jwtValidator.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Message{Message: "Invalid token"})
			c.Abort()
			return
		}

		// luu claim vao trong context
		c.Set("claims", tokenClaims)
		c.Next()
	}

}

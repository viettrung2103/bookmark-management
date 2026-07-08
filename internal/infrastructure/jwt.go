package infrastructure

import (
	"github.com/viettrung2103/bookmark-management/pkg/common"
	"github.com/viettrung2103/bookmark-management/pkg/jwtutils"
)

// CreateJWTProvider creates a new JWT generator and validator
func CreateJWTProvider() (jwtutils.JWTGenerator, jwtutils.JWTValidator) {
	jwtGen, err := jwtutils.NewJWTGenerator("./private.pem")
	common.HandleError(err)
	jwtVal, err := jwtutils.NewJWTValidator("./public.pem")
	common.HandleError(err)

	return jwtGen, jwtVal
}

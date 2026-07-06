package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/viettrung2103/bookmark-management/internal/app/service/user"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	"github.com/viettrung2103/bookmark-management/pkg/requestutils"
	"github.com/viettrung2103/bookmark-management/pkg/response"
)

type loginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=8"`
}

// Login Authentication endpoint
// @Summary return a jwt token if the input is correct
// @Description Return a jwt toke if the input is correct
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param input body loginInput tru "Input required"
// @Success 200 {object} object{token=string, message=string} "Success"
// @Router /v1/users/login [post]
func (u *userHandler) Login(c *gin.Context) {
	// get input
	//input := &loginInput{}
	//if err := c.ShouldBindJSON(input); err != nil {
	//	c.JSON(http.StatusBadRequest, response.InputFieldError(err))
	//	return
	//}
	input, err := requestutils.BindInputFromRequest[loginInput](c)
	if err != nil {
		return
	}

	// call service
	tokenStr, err := u.service.Login(c, input.Username, input.Password)
	switch {
	case errors.Is(err, user.ErrInvalidCreditials):
		c.JSON(http.StatusBadRequest, response.Message{Message: "Invalid username or password"})
		return
	case errors.Is(err, dbutils.ErrRecordNotFound):
		c.JSON(http.StatusBadRequest, response.Message{Message: "Invalid username or password"})
		return
	case err == nil:
	default:
		log.Err(err).Msg("Failed to login")
		c.JSON(http.StatusInternalServerError, response.Message{Message: "Internal server error"})
		return

	}

	// return response with token
	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

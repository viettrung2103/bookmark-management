package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	"github.com/viettrung2103/bookmark-management/pkg/requestutils"
	"github.com/viettrung2103/bookmark-management/pkg/response"
)

type registerInput struct {
	DisplayName string `json:"display_name" binding:"required,gt=0"`
	Username    string `json:"username" binding:"required,gt=0"`
	Password    string `json:"password" binding:"required,gt=8"`
	Email       string `json:"email" binding:"required,email"`
}

type createUser struct {
	Data    *model.User `json:"data"`
	Message string      `json:"message"`
}

// Register handles user registration
//
//	@Summary Create a netype createUser struct {
//		Data    *model.User `json:"data"`
//		Message string      `json:"message"`
//	}w user
//
// @Description Create a new user
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param input body registerInput true "User registration input"
// @Success 200 {object} object{data=model.User,message=string} "Success"
// @Router /v1/users/register [post]
func (h *userHandler) Register(c *gin.Context) {

	input, err := requestutils.BindInputFromRequest[registerInput](c)
	if err != nil {
		return
	}

	user, err := h.service.CreateUser(c, input.DisplayName, input.Username, input.Password, input.Email)

	switch {
	case errors.Is(err, dbutils.ErrDuplication):
		c.JSON(http.StatusBadRequest, response.Message{
			Message: "User or Email already exist",
		})
		return
	case err == nil:
		//return
	default:
		log.Err(err).Msg("Failed to create user")
		c.JSON(http.StatusInternalServerError, response.Message{
			Message: "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "Register an user successfully",
	})
}

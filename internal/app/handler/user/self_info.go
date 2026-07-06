package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/viettrung2103/bookmark-management/pkg/dbutils"
	"github.com/viettrung2103/bookmark-management/pkg/requestutils"
	"github.com/viettrung2103/bookmark-management/pkg/response"
)

// SelfInfo get your current user information
// @Summary get your current user information
// @Descripton get your current user information
// @Tags user
// @Security BearerAuth
// @Accept applicatin/json
// @Produce application/json
// @Success 200 {object} object{data=model.User} "Success"
// @Router /v1/self/info [get]
func (h *userHandler) SelfInfo(c *gin.Context) {

	userId, err := requestutils.GetUserIDFromRequest(c)
	if err != nil {
		return
	}

	user, err := h.service.SelfInfo(c, userId)
	switch {
	case errors.Is(err, dbutils.ErrRecordNotFound):
		log.Error().Err(err).Str("userId", userId).Msg("User Not Found")
		c.JSON(http.StatusBadRequest, &response.Message{
			Message: "User Not Found",
		})
		return
	case errors.Is(err, nil):
		break
	default:
		log.Err(err).Str("userId", userId).Msg("failed to fetch user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, user)
}

type editInput struct {
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
}

// EditSelfInfo edit your current user info
// @Summary edit your current info
// @Description edit your current user info
// @Tags user
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Param input body editInput true "Input required"
// @Success 200 {object} object{message=string} "Success"
// @Router /v1/self/info [put]
func (h *userHandler) EditSelfInfo(c *gin.Context) {
	// GETURL the user id from jwt token
	uid, err := requestutils.GetUserIDFromRequest(c)
	if err != nil {
		return
	}

	// getting input from request and valid
	input := &editInput{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	err = h.service.EditInfoByID(c, uid, input.DisplayName, input.Email)
	switch {
	case errors.Is(err, dbutils.ErrDuplication):
		c.JSON(http.StatusBadRequest, &response.Message{
			Message: err.Error(),
		})
		return
	case errors.Is(err, nil):
		break
	default:
		log.Err(err).Str("userID", uid).Msg("failed to update user info")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Edit current user successfully!",
	})

}

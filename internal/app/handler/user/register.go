package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerInput struct {
	DisplayName string `json:"display_name" binding:"required,gt=0"`
	Username    string `json:"username" binding:"required,gt=0"`
	Password    string `json:"password" binding:"required,gt=8"`
	Email       string `json:"email" binding:"required,email"`
}

// Register handles user registration
// @Summary Create a new user
// @Description Create a new user
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param input body registerInput true "User registration input"
// @Success 200 {object} object{data=model.User,message=string} "Success"
// @Router /v1/users/register [post]
func (h *userHandler) Register(c *gin.Context) {
	input := &registerInput{}
	if err := c.ShouldBindJSON(input); err != nil {
		//c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	user, err := h.service.CreateUser(c, input.DisplayName, input.Username, input.Password, input.Email)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

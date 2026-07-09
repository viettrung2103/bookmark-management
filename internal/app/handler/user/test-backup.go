package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/viettrung2103/bookmark-management/pkg/response"
)

type testInput struct {
	ID    string `uri:"id" validate:"required"`
	UID   string `uri:"uid"`
	Count int    `header:"count" validate:"required"`
	Name  string `form:"name" validate:"required"`
}

func Test(c *gin.Context) {
	input := &testInput{}
	if err := c.ShouldBindHeader(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}
	if err := c.ShouldBindUri(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	if err := c.ShouldBindQuery(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	inputValidator := validator.New(validator.WithRequiredStructEnabled())
	if err := inputValidator.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	println("this is input", input)
	println("this is UID", input.UID)
	println("this is ID", input.ID)
	println("this is header", input.Count)
	println("this is param name", input.Name)

	fmt.Sprintf("ID: %s, UID: %s, Count: %d, Name: %s", input.ID, input.UID, input.Count, input.Name)

	c.JSON(http.StatusOK, fmt.Sprintf("ID: %s, UID: %s, Count:%d, Name: %s", input.ID, input.UID, input.Count, input.Name))
}

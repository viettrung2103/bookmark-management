package response

import "github.com/gin-gonic/gin"

type Message struct {
	Message string `json:"message"`
}

var InternalErrResponse = Message{
	Message: "Internal server error",
}

func InputFieldError(err error) gin.H {
	return gin.H{
		"error":   "the input is invalid",
		"message": err.Error(),
	}
}

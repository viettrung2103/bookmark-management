package response

import "github.com/gin-gonic/gin"

// Message represents a response message
type Message struct {
	Message string `json:"message"`
}

var InternalErrResponse = Message{
	Message: "Internal server error",
}

// InputFieldError returns a map with error and message fields
func InputFieldError(err error) gin.H {
	return gin.H{
		"error":   "the input is invalid",
		"message": err.Error(),
	}
}

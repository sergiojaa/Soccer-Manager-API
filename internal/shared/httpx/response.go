package httpx

import "github.com/gin-gonic/gin"

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(c *gin.Context, status int, data interface{}, message string) {
	response := gin.H{
		"data": data,
	}

	if message != "" {
		response["message"] = message
	}

	c.JSON(status, response)
}

func Error(c *gin.Context, status int, code string, message string) {
	c.JSON(status, gin.H{
		"error": ErrorPayload{
			Code:    code,
			Message: message,
		},
	})
}

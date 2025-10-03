package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseWithData is a helper function to return a successful response with data.
func ResponseWithData(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    data,
	})
}

func ResponseWithToken(c *gin.Context, message string, email, token, expiryDate, tokenType string) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data": gin.H{
			"email": email,
			"tokens": gin.H{
				"access_token": token,
				"token_type":   tokenType,
				"expiry_date":  expiryDate,
			},
		},
	})
}

// ResponseWithMessage is a helper function to return a successful response with just a message.
func ResponseWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

// ErrorMessage is a helper function to return an error response with a message and status code.
func ErrorMessage(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{
		"message": message,
	})
}

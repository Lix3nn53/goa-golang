package error

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// Generic Errors
var (
	// ErrNotFound error will be returned when not found a model
	ErrNotFound = errors.New("not found")
	// InvalidPaymentMethod error will be returned when payment method is not available
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
)

func Respond(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}

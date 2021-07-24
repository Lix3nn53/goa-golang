package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type testMiddleware struct{}

//TestMiddlewareInterface ...
type TestMiddlewareInterface interface {
	Handler() gin.HandlerFunc
}

//NewTestMiddleware ...
func NewTestMiddleware() TestMiddlewareInterface {
	return &testMiddleware{}
}

//Handler ...
func (cm testMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(("TEST MIDDLEWARE"))

		c.Next()
	}
}

package middleware

import (
	"github.com/gin-gonic/gin"

	"goa-golang/internal/logger"
)

type testMiddleware struct {
	logger logger.Logger
}

//TestMiddlewareInterface ...
type TestMiddlewareInterface interface {
	Handler() gin.HandlerFunc
}

//NewTestMiddleware ...
func NewTestMiddleware(logger logger.Logger) TestMiddlewareInterface {
	return &testMiddleware{
		logger,
	}
}

//Handler ...
func (cm testMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		cm.logger.Infof("Test Middleware: %s", c.ClientIP())

		c.Next()
	}
}

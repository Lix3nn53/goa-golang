package middleware

import (
	"github.com/gin-gonic/gin"
)

type corsMiddleware struct{}

//CorsMiddlewareInterface ...
type CorsMiddlewareInterface interface {
	Handler() gin.HandlerFunc
}

//NewCorsMiddleware ...
func NewCorsMiddleware() CorsMiddlewareInterface {
	return &corsMiddleware{}
}

//Handler ...
func (cm corsMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

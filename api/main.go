package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/go-openapi/runtime/middleware"
)

func DummyMiddleware(c *gin.Context) {
	fmt.Println("Im a dummy!")

	// Pass on to the next-in-chain
	c.Next()
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/test", DummyMiddleware, func(c *gin.Context) {
		c.String(200, "test successful")
	})

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yml"}
	sh := middleware.Redoc(opts, nil)

	r.GET("/docs", gin.WrapH(sh))

	r.StaticFile("/swagger.yml", "./public/swagger.yml")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")

	return r
}

func main() {
	fmt.Println("Hello World!")

	r := setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

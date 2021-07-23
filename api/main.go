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

func Ping(c *gin.Context) {
	c.String(200, "pong")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Simple group: v1
	v1 := r.Group("/v1")

	v1.Use(DummyMiddleware)
	{
		v1.StaticFile("/favicon.ico", "./public/favicon.ico")
		v1.GET("/ping", Ping)

		v1.GET("/test", DummyMiddleware, func(c *gin.Context) {
			c.String(200, "test successful")
		})

		// handler for documentation
		opts := middleware.RedocOpts{BasePath: "/v1", SpecURL: "/v1/swagger.yml"}
		sh := middleware.Redoc(opts, nil)
		v1.GET("/docs", gin.WrapH(sh))
		v1.StaticFile("/swagger.yml", "./public/swagger.yml")
	}

	return r
}

func main() {
	fmt.Println("Hello World!")

	r := setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

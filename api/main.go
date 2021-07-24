package main

import (
	"fmt"

	"goa-golang/internal/dic"
	"goa-golang/internal/logger"
	"goa-golang/internal/route"

	"flag"
	"os"

	"github.com/joho/godotenv"
)

// func DummyMiddleware(c *gin.Context) {
// 	fmt.Println("Im a dummy!")

// 	// Pass on to the next-in-chain
// 	c.Next()
// }

// func Ping(c *gin.Context) {
// 	c.String(200, "pong")
// }

// func setupRouter() *gin.Engine {
// 	r := gin.Default()

// 	// Simple group: v1
// 	v1 := r.Group("/v1")

// 	v1.Use(DummyMiddleware)
// 	{
// 		v1.StaticFile("/favicon.ico", "./public/favicon.ico")
// 		v1.GET("/ping", Ping)

// 		v1.GET("/test", DummyMiddleware, func(c *gin.Context) {
// 			c.String(200, "test successful")
// 		})

// 		// handler for documentation
// 		opts := middleware.RedocOpts{BasePath: "/v1", SpecURL: "/v1/swagger.yml"}
// 		sh := middleware.Redoc(opts, nil)
// 		v1.GET("/docs", gin.WrapH(sh))
// 		v1.StaticFile("/swagger.yml", "./public/swagger.yml")
// 	}

// 	return r
// }

var config string

func main() {
	fmt.Println("Hello World!")

	flag.StringVar(&config, "env", "dev.env", "Environment name")
	flag.Parse()

	logger := logger.NewAPILogger()
	logger.InitLogger()

	if err := godotenv.Load(config); err != nil {
		logger.Fatalf(err.Error())
	}
	container := dic.InitContainer()
	router := route.Setup(container, logger)
	router.Run(":" + os.Getenv("APP_PORT"))
}

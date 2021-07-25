package main

import (
	"goa-golang/internal/dic"
	"goa-golang/internal/logger"
	"goa-golang/internal/route"

	"flag"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var config string

func setupRouter() (*gin.Engine, logger.Logger) {
	flag.StringVar(&config, "env", "dev.env", "Environment name")
	flag.Parse()

	logger := logger.NewAPILogger()
	logger.InitLogger()

	logger.Info("Hello World!")

	if err := godotenv.Load(config); err != nil {
		logger.Fatalf(err.Error())
	}
	container := dic.InitContainer(logger)
	router := route.Setup(container, logger)

	return router, logger
}

func main() {
	router, _ := setupRouter()

	router.Run(":" + os.Getenv("APP_PORT"))
}

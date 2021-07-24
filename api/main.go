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
	container := dic.InitContainer(logger)
	router := route.Setup(container, logger)
	router.Run(":" + os.Getenv("APP_PORT"))
}

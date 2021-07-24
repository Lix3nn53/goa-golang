package config

import (
	"os"
)

var Conf *Config

/**
 * Config
 *
 */
type Config struct {
	Port      string
	BaseUrl   string
	Db        Database
	SecretKey string
	// Email     Email
	Brand struct {
		ProjectName   string
		ProjectUrl    string
		ProjectApiUrl string
	}
}

/**
 * Configurations
 *
 */
func Configurations() {
	port := os.Getenv("API_PORT")
	Conf = &Config{
		Port:      port,
		BaseUrl:   os.Getenv("BASE_URL") + ":" + port,
		SecretKey: os.Getenv("JWT_SECRET_KEY"),
		// Email:     GetEmailConfig(),
		Brand: struct {
			ProjectName   string
			ProjectUrl    string
			ProjectApiUrl string
		}{ProjectName: os.Getenv("PROJECT_NAME"), ProjectUrl: os.Getenv("PROJECT_URL"), ProjectApiUrl: os.Getenv("PROJECT_API_URL")},
	}
}

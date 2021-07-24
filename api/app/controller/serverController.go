package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//ServerControllerInterface define the user controller interface methods
type ServerControllerInterface interface {
	Ping(c *gin.Context)
	Version(c *gin.Context)
}

// serverController handles communication with the user service
type serverController struct{}

// NewServerController implements the user controller interface.
func NewServerController() ServerControllerInterface {
	return &serverController{}
}

func (uc *serverController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (uc *serverController) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": os.Getenv("VERSION")})
}

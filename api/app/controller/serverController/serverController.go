package serverController

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

type PingResponse struct {
	Message string `json:"message"`
}

func (uc *serverController) Ping(c *gin.Context) {
	response := PingResponse{Message: "pong"}
	c.JSON(http.StatusOK, response)
}

type VersionResponse struct {
	Version string `json:"version"`
}

func (uc *serverController) Version(c *gin.Context) {
	response := VersionResponse{Version: os.Getenv("VERSION")}
	c.JSON(http.StatusOK, response)
}

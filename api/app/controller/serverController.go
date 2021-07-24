package controller

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"os"
)

type ServerController struct{}

// Ping godoc
// @Tags Server
// @Success 200 {object} viewModels.Message{}
// @Failure 500
// @Router /status/ping [get]
func (ServerController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// Version godoc
// @Tags Server
// @Success 200 {object} viewModels.Message{}
// @Failure 500
// @Router /status/version [get]
func (ServerController) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": os.Getenv("VERSION")})
}

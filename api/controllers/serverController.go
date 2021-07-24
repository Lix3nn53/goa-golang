package controllers

import (
	"github.com/gin-gonic/gin"

	"goa-golang/viewModels"
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
	c.JSON(http.StatusOK, viewModels.MResponse("pong"))
}

// Version godoc
// @Tags Server
// @Success 200 {object} viewModels.Message{}
// @Failure 500
// @Router /status/version [get]
func (ServerController) Version(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"version": os.Getenv("VERSION")})
}

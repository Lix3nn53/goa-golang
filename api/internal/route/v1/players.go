package v1

import (
	"goa-golang/app/controller/authController"
	"goa-golang/app/controller/playerController"

	"github.com/gin-gonic/gin"
)

func SetupPlayerRoute(players *gin.RouterGroup, pc playerController.PlayerControllerInterface, ac authController.AuthControllerInterface) *gin.RouterGroup {
	players.Use(ac.AuthMiddleware())

	players.GET("/info", pc.Info)

	return players
}

package playerController

import (
	"errors"
	appError "goa-golang/app/error"
	"goa-golang/app/service/playerService"
	"goa-golang/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

//PlayerControllerInterface define the player controller interface methods
type PlayerControllerInterface interface {
	Info(c *gin.Context)
}

// PlayerController handles communication with the player service
type PlayerController struct {
	service playerService.PlayerServiceInterface
	logger  logger.Logger
}

// NewPlayerController implements the player controller interface.
func NewPlayerController(service playerService.PlayerServiceInterface, logger logger.Logger) PlayerControllerInterface {
	return &PlayerController{
		service,
		logger,
	}
}

// Find implements the method to handle the service to find a player by the primary key
func (uc *PlayerController) Info(c *gin.Context) {
	playerUUID, exists := c.Get("playerUUID")

	if !exists {
		appError.Respond(c, http.StatusForbidden, errors.New("no playerUUID"))
		return
	}

	uuid := playerUUID.(string)

	player, err := uc.service.FindByID(uuid)
	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, player)
}

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
	tokenId, exists := c.Get("tokenID")
	if !exists {
		appError.Respond(c, http.StatusForbidden, errors.New("no id"))
		return
	}

	tokenIdField, exists := c.Get("tokenIDField")
	if !exists {
		appError.Respond(c, http.StatusForbidden, errors.New("no idField"))
		return
	}

	id := tokenId.(string)
	idField := tokenIdField.(string)

	if idField != "uuid" {
		appError.Respond(c, http.StatusForbidden, errors.New("not implemented"))
		return
	}

	player, err := uc.service.FindByID(id)
	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, player)
}

package characterController

import (
	"errors"
	appError "goa-golang/app/error"
	"goa-golang/app/service/characterService"
	"goa-golang/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

//CharacterControllerInterface define the character controller interface methods
type CharacterControllerInterface interface {
	Info(c *gin.Context)
}

// CharacterController handles communication with the character service
type CharacterController struct {
	service characterService.CharacterServiceInterface
	logger  logger.Logger
}

// NewCharacterController implements the character controller interface.
func NewCharacterController(service characterService.CharacterServiceInterface, logger logger.Logger) CharacterControllerInterface {
	return &CharacterController{
		service,
		logger,
	}
}

// Find implements the method to handle the service to find a character by the primary key
func (uc *CharacterController) Info(c *gin.Context) {
	userUUID, exists := c.Get("userUUID")

	if !exists {
		appError.Respond(c, http.StatusForbidden, errors.New("no userUUID"))
		return
	}

	uuid := userUUID.(string)

	characters, err := uc.service.FindByID(uuid)
	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, characters)
}

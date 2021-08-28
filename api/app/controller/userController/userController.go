package userController

import (
	"errors"
	appError "goa-golang/app/error"
	"goa-golang/app/service/userService"
	"goa-golang/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

//UserControllerInterface define the user controller interface methods
type UserControllerInterface interface {
	Info(c *gin.Context)
}

// UserController handles communication with the user service
type UserController struct {
	service userService.UserServiceInterface
	logger  logger.Logger
}

// NewUserController implements the user controller interface.
func NewUserController(service userService.UserServiceInterface, logger logger.Logger) UserControllerInterface {
	return &UserController{
		service,
		logger,
	}
}

// Find implements the method to handle the service to find a user by the primary key
func (uc *UserController) Info(c *gin.Context) {
	tokenId, exists := c.Get("token_id")
	if !exists {
		appError.Respond(c, http.StatusForbidden, errors.New("no id"))
		return
	}

	tokenIdField, exists := c.Get("token_idField")
	if !exists {
		appError.Respond(c, http.StatusForbidden, errors.New("no idField"))
		return
	}

	id := tokenId.(string)
	idField := tokenIdField.(string)

	user, err := uc.service.FindByID(id, idField)
	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

package userController

import (
	errorNotFound "goa-golang/app/error"
	"goa-golang/app/service/authService"
	"goa-golang/internal/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//UserControllerInterface define the user controller interface methods
type AuthControllerInterface interface {
	GoogleOauth2(c *gin.Context)
}

// UserController handles communication with the user service
type AuthController struct {
	service authService.AuthServiceInterface
	logger  logger.Logger
}

// NewUserController implements the user controller interface.
func NewUserController(service authService.AuthServiceInterface, logger logger.Logger) AuthControllerInterface {
	return &AuthController{
		service,
		logger,
	}
}

// Find implements the method to handle the service to find a user by the primary key
func (uc *AuthController) GoogleOauth2(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := uc.service.GoogleOauth2(id)
	if err != nil {
		uc.logger.Error(err.Error())
		c.Status(errorNotFound.ParseError(err))
		return
	}
	c.JSON(http.StatusOK, user)
}

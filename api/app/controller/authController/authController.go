package authController

import (
	errorNotFound "goa-golang/app/error"
	"goa-golang/app/service/authService"
	"goa-golang/internal/logger"
	"net/http"

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
func NewAuthController(service authService.AuthServiceInterface, logger logger.Logger) AuthControllerInterface {
	return &AuthController{
		service,
		logger,
	}
}

// Find implements the method to handle the service to find a user by the primary key
func (uc *AuthController) GoogleOauth2(c *gin.Context) {
	code := c.Query("code")

	user, err := uc.service.GoogleOauth2(code)
	if err != nil {
		uc.logger.Error(err.Error())
		c.Status(errorNotFound.ParseError(err))
		return
	}
	c.JSON(http.StatusOK, user)
}

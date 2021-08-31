package mojangController

import (
	appError "goa-golang/app/error"
	"goa-golang/app/service/mojangService"
	"goa-golang/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

//MojangControllerInterface define the mojang controller interface methods
type MojangControllerInterface interface {
	Profiles(c *gin.Context)
}

// MojangController handles communication with the mojang service
type MojangController struct {
	service mojangService.MojangServiceInterface
	logger  logger.Logger
}

// NewMojangController implements the mojang controller interface.
func NewMojangController(service mojangService.MojangServiceInterface, logger logger.Logger) MojangControllerInterface {
	return &MojangController{
		service,
		logger,
	}
}

type RequestProfiles struct {
	UUID string `json:"uuid" binding:"required"`
}

// Find implements the method to handle the service to find a mojang by the primary key
func (uc *MojangController) Profiles(c *gin.Context) {
	var req RequestProfiles
	err := c.BindJSON(&req)
	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusBadRequest, err)
		return
	}

	mojang, err := uc.service.Profiles(req.UUID)
	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, mojang)
}

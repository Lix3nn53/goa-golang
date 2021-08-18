package userController

import (
	appError "goa-golang/app/error"
	"goa-golang/app/model/userModel"
	"goa-golang/app/service/userService"
	"goa-golang/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

//UserControllerInterface define the user controller interface methods
type UserControllerInterface interface {
	Find(c *gin.Context)
	Destroy(c *gin.Context)
	Update(c *gin.Context)
	Store(c *gin.Context)
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
func (uc *UserController) Find(c *gin.Context) {

	uuid := c.Param("uuid")

	user, err := uc.service.FindByID(uuid)
	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// Destroy implements the method to validate the params to store a  new user and handle the service
func (uc *UserController) Destroy(c *gin.Context) {
	uuid := c.Param("uuid")

	err := uc.service.RemoveByID(uuid)

	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusOK)
}

// Update implements the method to validate teh params to update a user and handle the service
func (uc *UserController) Update(c *gin.Context) {
	uuid := c.Param("id")

	var user userModel.UpdateUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = uc.service.UpdateByID(uuid, user)

	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusOK)
}

// Store implements the method to validate the params to store a  new user and handle the service
func (uc *UserController) Store(c *gin.Context) {

	var rq userModel.CreateUser

	if err := c.ShouldBindJSON(&rq); err != nil {
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	validate := validator.New()
	err := validate.Struct(rq)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := uc.service.Store("1", rq)

	if err != nil {
		uc.logger.Error(err.Error())
		appError.Respond(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, *user)
}

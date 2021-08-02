package billingController

import (
	errorNotFound "goa-golang/app/error"
	"goa-golang/app/model/billingModel"
	"goa-golang/app/service/billingService"
	"goa-golang/app/service/userService"
	"goa-golang/internal/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

//BillingControllerInterface define the services controller interface methods
type BillingControllerInterface interface {
	AddCustomer(c *gin.Context)
}

// billingController handles communication with the external service
type BillingController struct {
	service  billingService.BillingServiceInterface
	uservice userService.UserServiceInterface
	logger   logger.Logger
}

// NewBillingController implements the user controller interface.
func NewBillingController(service billingService.BillingServiceInterface, uservice userService.UserServiceInterface, logger logger.Logger) BillingControllerInterface {
	return &BillingController{
		service,
		uservice,
		logger,
	}
}

// Store implements the method to validate the params to store a  new payment method and handle the service
func (bc *BillingController) AddCustomer(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := bc.uservice.FindByID(id)
	if err != nil {
		bc.logger.Error(err.Error())
		c.Status(errorNotFound.ParseError(err))
		return
	}

	var rq billingModel.CreateCustomer

	if err := c.ShouldBindJSON(&rq); err != nil {
		bc.logger.Error(err.Error())
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	validate := validator.New()

	err = validate.Struct(rq)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := bc.service.GetPaymentAdapter(rq)
	if err != nil {
		bc.logger.Error(err.Error())
		c.Status(errorNotFound.ParseError(err))
		return
	}
	err = bc.service.AddBilling(*user, *p)

	if err != nil {
		bc.logger.Error(err.Error())
		c.Status(errorNotFound.ParseError(err))
		return
	}

	c.Status(http.StatusOK)
}

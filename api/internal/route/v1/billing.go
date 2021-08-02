package v1

import (
	"goa-golang/app/controller/billingController"

	"github.com/gin-gonic/gin"
)

func SetupBillingRoute(userBilling *gin.RouterGroup, bc billingController.BillingControllerInterface) *gin.RouterGroup {
	userBilling.POST("", bc.AddCustomer)

	return userBilling
}

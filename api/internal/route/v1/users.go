package v1

import (
	"goa-golang/app/controller"

	"github.com/gin-gonic/gin"
)

func SetupUserRoute(users *gin.RouterGroup, uc controller.UserControllerInterface) *gin.RouterGroup {
	users.POST("", uc.Store)

	user := users.Group(":id")
	{
		user.GET("", uc.Find)
		user.DELETE("", uc.Destroy)
		user.PUT("", uc.Update)
		// userBilling := user.Group("/billing")
		// {
		// 	userBilling.POST("", billc.AddCustomer)
		// }
	}

	return users
}

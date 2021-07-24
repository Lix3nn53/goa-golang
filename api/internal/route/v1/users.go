package v1

import (
	"goa-golang/app/controller"
	"goa-golang/app/service"
	"goa-golang/internal/dic"
	"goa-golang/internal/logger"

	"github.com/gin-gonic/gin"

	"github.com/sarulabs/dingo/generation/di"
)

func SetupUserRoute(users *gin.RouterGroup, container di.Container, logger logger.Logger) *gin.RouterGroup {
	uc := controller.NewUserController(container.Get(dic.UserService).(service.UserServiceInterface), logger)

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

package v1

import (
	"goa-golang/app/controller/userController"

	"github.com/gin-gonic/gin"
)

func SetupUserRoute(users *gin.RouterGroup, uc userController.UserControllerInterface) *gin.RouterGroup {
	users.POST("", uc.Store)

	user := users.Group(":id")
	{
		user.GET("", uc.Find)
		user.DELETE("", uc.Destroy)
		user.PUT("", uc.Update)
	}

	return users
}

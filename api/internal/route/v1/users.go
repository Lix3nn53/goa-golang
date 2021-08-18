package v1

import (
	"goa-golang/app/controller/userController"

	"github.com/gin-gonic/gin"
)

func SetupUserRoute(users *gin.RouterGroup, uc userController.UserControllerInterface) *gin.RouterGroup {
	users.GET("/info", uc.Info)

	return users
}

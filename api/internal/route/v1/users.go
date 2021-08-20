package v1

import (
	"goa-golang/app/controller/authController"
	"goa-golang/app/controller/userController"

	"github.com/gin-gonic/gin"
)

func SetupUserRoute(users *gin.RouterGroup, uc userController.UserControllerInterface, ac authController.AuthControllerInterface) *gin.RouterGroup {
	users.Use(ac.AuthMiddleware())

	users.GET("/info", uc.Info)

	return users
}

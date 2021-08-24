package v1

import (
	"goa-golang/app/controller/authController"
	"goa-golang/app/controller/characterController"

	"github.com/gin-gonic/gin"
)

func SetupCharacterRoute(characters *gin.RouterGroup, c characterController.CharacterControllerInterface, ac authController.AuthControllerInterface) *gin.RouterGroup {
	characters.Use(ac.AuthMiddleware())

	characters.GET("/info", c.Info)

	return characters
}

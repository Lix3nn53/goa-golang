package v1

import (
	"goa-golang/app/controller/authController"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoute(auth *gin.RouterGroup, ac authController.AuthControllerInterface) *gin.RouterGroup {
	auth.GET("google", ac.GoogleOauth2)

	return auth
}

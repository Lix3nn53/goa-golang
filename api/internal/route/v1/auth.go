package v1

import (
	"goa-golang/app/controller/authController"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoute(auth *gin.RouterGroup, c authController.AuthControllerInterface) *gin.RouterGroup {
	auth.GET("refresh_token", c.RefreshAccessToken)
	auth.GET("logout", c.Logout)
	auth.GET("google", c.GoogleOauth2)

	return auth
}

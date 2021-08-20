package v1

import (
	"goa-golang/app/controller/authController"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoute(auth *gin.RouterGroup, ac authController.AuthControllerInterface) *gin.RouterGroup {
	auth.GET("refresh_access_token", ac.RefreshAccessToken)
	auth.GET("logout", ac.Logout)
	auth.GET("google", ac.GoogleOauth2)

	return auth
}

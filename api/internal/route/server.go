package route

import (
	"github.com/gin-gonic/gin"

	"goa-golang/app/controller"
)

func SetupServerRoute(r *gin.Engine) *gin.Engine {
	//server
	r.StaticFile("/favicon.ico", "../public/favicon.ico")

	sc := controller.NewServerController()

	r.GET("/ping", sc.Ping)
	r.GET("/version", sc.Version)

	return r
}

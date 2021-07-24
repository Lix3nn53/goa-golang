package route

import (
	"github.com/gin-gonic/gin"

	"goa-golang/app/controllers"
)

func SetupServerRoute(r *gin.Engine) *gin.Engine {
	//server
	r.StaticFile("/favicon.ico", "../public/favicon.ico")

	status := r.Group("/status")
	status.GET("/status/ping", controllers.ServerController{}.Ping)
	status.GET("/status/version", controllers.ServerController{}.Version)

	return r
}

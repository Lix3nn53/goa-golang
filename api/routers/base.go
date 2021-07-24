package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"

	"goa-golang/app"
	"goa-golang/controllers"
)

func Route(r *gin.Engine) *gin.Engine {
	//server
	r.StaticFile("/favicon.ico", "../public/favicon.ico")
	r.GET("/status/ping", controllers.ServerController{}.Ping)
	r.GET("/status/version", controllers.ServerController{}.Version)

	// Simple group: v1
	v1 := r.Group("/v1")

	v1.Use(app.Application.Container.GetAuthMiddleware().AuthMiddleware)
	{
		// handler for documentation
		opts := middleware.RedocOpts{BasePath: "/v1", SpecURL: "/v1/swagger.yml"}
		sh := middleware.Redoc(opts, nil)
		v1.GET("/docs", gin.WrapH(sh))
		v1.StaticFile("/swagger.yml", "./public/swagger.yml")
	}

	return r
}

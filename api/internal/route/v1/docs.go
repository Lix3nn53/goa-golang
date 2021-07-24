package v1

import (
	"goa-golang/internal/dic"
	"goa-golang/internal/middleware"

	"github.com/gin-gonic/gin"
	openapi "github.com/go-openapi/runtime/middleware"
	"github.com/sarulabs/dingo/generation/di"
)

func SetupDocsRoute(v1 *gin.RouterGroup, container di.Container) *gin.RouterGroup {
	testMiddleware := container.Get(dic.TestMiddleware).(middleware.TestMiddlewareInterface)

	v1.Use(testMiddleware.Handler())
	{
		// handler for documentation
		opts := openapi.RedocOpts{BasePath: "/v1", SpecURL: "/v1/swagger.yml"}
		sh := openapi.Redoc(opts, nil)

		v1.GET("/docs", gin.WrapH(sh))
		v1.StaticFile("/swagger.yml", "./public/swagger.yml")
	}

	return v1
}

package v1

import (
	"goa-golang/internal/middleware"

	"github.com/gin-gonic/gin"
	openapi "github.com/go-openapi/runtime/middleware"
)

func SetupDocsRoute(v1 *gin.RouterGroup, testMiddleware middleware.TestMiddlewareInterface) *gin.RouterGroup {
	v1.Use(testMiddleware.Handler())
	{
		// handler for documentation
		v1.StaticFile("/swagger.yml", "./public/docs/swagger.yml")

		opts := openapi.RedocOpts{BasePath: "/v1", SpecURL: "/v1/swagger.yml"}
		sh := openapi.Redoc(opts, nil)

		v1.GET("/docs", gin.WrapH(sh))
	}

	return v1
}

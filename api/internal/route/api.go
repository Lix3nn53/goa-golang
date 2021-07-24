package route

import (
	"fmt"
	routev1 "goa-golang/app/route/v1"
	"goa-golang/internal/dic"
	"goa-golang/internal/logger"
	"goa-golang/internal/middleware"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/dingo/generation/di"
)

// Setup returns initialized routes.
func Setup(container di.Container, logger logger.Logger) *gin.Engine {
	// ac := container.Get(dic.AuthController).(controller.AuthControllerInterface)

	gin.SetMode(os.Getenv("GIN_MODE"))

	r := gin.New()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	r.Use(gin.Recovery())

	// Middleware initialization
	corsMiddleware := container.Get(dic.CorsMiddleware).(middleware.CorsMiddlewareInterface)
	r.Use(corsMiddleware.Handler())

	// server Routes
	SetupServerRoute(r)

	// v1 Routes

	// uc := container.Get(dic.UserController).(controller.UserControllerInterface)
	// uc := controller.NewUserController(container.Get(dic.UserService).(service.UserServiceInterface), logger)

	// billc := controller.NewBillingController(container.Get(dic.BillingService).(service.BillingServiceInterface), container.Get(dic.UserService).(service.UserServiceInterface), logger)

	v1 := r.Group("/api")

	routev1.SetupDocsRoute(v1)
	routev1.SetupUserRoute(v1, container, logger)

	return r
}

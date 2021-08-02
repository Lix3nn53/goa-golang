package route

import (
	"fmt"
	"goa-golang/internal/dic"
	"goa-golang/internal/logger"
	"goa-golang/internal/middleware"
	routev1 "goa-golang/internal/route/v1"
	"goa-golang/internal/storage"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Setup returns initialized routes.
func Setup(db *storage.DbStore, dbCache *storage.DbCache, logger logger.Logger) *gin.Engine {
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
	corsMiddleware := middleware.NewCorsMiddleware()
	r.Use(corsMiddleware.Handler())

	// server Routes
	SetupServerRoute(r)

	// v1 Routes
	v1 := r.Group("/v1")
	{
		testMiddleware := middleware.NewTestMiddleware(logger)
		routev1.SetupDocsRoute(v1, testMiddleware)

		users := v1.Group("/users")
		{
			userController, userService := dic.InitUserController(db, logger)
			routev1.SetupUserRoute(users, userController)

			user := users.Group(":id")
			{
				billingController := dic.InitBillingController(db, userService, logger)
				userBilling := user.Group("/billing")
				{
					routev1.SetupBillingRoute(userBilling, billingController)
				}
			}
		}
	}

	return r
}

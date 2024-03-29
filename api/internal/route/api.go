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
		routev1.SetupDocsRoute(v1)

		// user
		userRepo := dic.InitUserRepository(db)

		// player
		playerRepo := dic.InitPlayerRepository(db)

		// auth
		authService := dic.InitAuthService(playerRepo, userRepo, logger)
		authCont := dic.InitAuthController(authService, logger)

		auth := v1.Group("/auth")
		{

			routev1.SetupAuthRoute(auth, authCont)
		}

		users := v1.Group("/users")
		{
			userService := dic.InitUserService(userRepo)
			userCont := dic.InitUserController(userService, logger)

			routev1.SetupUserRoute(users, userCont, authCont)
		}

		players := v1.Group("/players")
		{
			playerService := dic.InitPlayerService(playerRepo)
			playerCont := dic.InitPlayerController(playerService, logger)

			routev1.SetupPlayerRoute(players, playerCont, authCont)
		}

		characters := v1.Group("/characters")
		{
			characterRepo := dic.InitCharacterRepository(db)
			characterService := dic.InitCharacterService(characterRepo)
			characterCont := dic.InitCharacterController(characterService, logger)

			routev1.SetupPlayerRoute(characters, characterCont, authCont)
		}

		mojang := v1.Group("/mojang")
		{
			mojangService := dic.InitMojangService(logger)
			mojangCont := dic.InitMojangController(mojangService, logger)

			routev1.SetupMojangRoute(mojang, mojangCont)
		}
	}

	return r
}

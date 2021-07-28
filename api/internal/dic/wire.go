//+build wireinject

// The build tag makes sure the stub is not built in the final build.
package dic

import (
	"goa-golang/app/controller"
	"goa-golang/app/repository"
	"goa-golang/app/service"
	"goa-golang/internal/logger"
	"goa-golang/internal/storage"

	"github.com/google/wire"
)

func initUserController(db *storage.DbStore, logger logger.Logger) controller.UserControllerInterface {
	wire.Build(repository.NewUserRepository, service.NewUserService, controller.NewUserController)

	return &controller.UserController{}
}

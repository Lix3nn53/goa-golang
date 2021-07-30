//+build wireinject

// The build tag makes sure the stub is not built in the final build.

//lint:file-ignore U1000 Ignore all unused code
package dic

import (
	userController "goa-golang/app/controller/user"
	"goa-golang/app/repository"
	"goa-golang/app/service"
	"goa-golang/internal/logger"
	"goa-golang/internal/storage"

	"github.com/google/wire"
)

func initUserController(db *storage.DbStore, logger logger.Logger) userController.UserControllerInterface {
	wire.Build(repository.NewUserRepository, service.NewUserService, userController.NewUserController)

	return &userController.UserController{}
}

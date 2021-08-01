//+build wireinject
//
// The build tag makes sure the stub is not built in the final build.
//
//lint:file-ignore U1000 Ignore all unused code

package dic

import (
	"goa-golang/app/controller/userController"
	"goa-golang/app/repository/userRepository"
	"goa-golang/app/service/userService"
	"goa-golang/internal/logger"
	"goa-golang/internal/storage"

	"github.com/google/wire"
)

func initUserController(db *storage.DbStore, logger logger.Logger) userController.UserControllerInterface {
	wire.Build(userRepository.NewUserRepository, userService.NewUserService, userController.NewUserController)

	return &userController.UserController{}
}

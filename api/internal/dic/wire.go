//+build wireinject
//
// The build tag makes sure the stub is not built in the final build.
//
//lint:file-ignore U1000 Ignore all unused code

package dic

import (
	"goa-golang/app/controller/userController"
	"goa-golang/app/repository/billingRepository"
	"goa-golang/app/repository/userRepository"
	"goa-golang/app/service/billingService"
	"goa-golang/app/service/userService"
	"goa-golang/internal/logger"
	"goa-golang/internal/storage"

	"github.com/google/wire"
)

func initUserService(db *storage.DbStore) userService.UserServiceInterface {
	wire.Build(userRepository.NewUserRepository, userService.NewUserService)

	return &userService.UserService{}
}

func initUserController(us userService.UserServiceInterface, logger logger.Logger) userController.UserControllerInterface {
	wire.Build(userController.NewUserController)

	return &userController.UserController{}
}

func initBillingService(db *storage.DbStore) billingService.BillingServiceInterface {
	wire.Build(billingRepository.NewBillingRepository, billingService.NewBillingService)

	return &billingService.BillingService{}
}

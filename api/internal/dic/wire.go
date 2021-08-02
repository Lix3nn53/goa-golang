//+build wireinject
//
// The build tag makes sure the stub is not built in the final build.
//
//lint:file-ignore U1000 Ignore all unused code

package dic

import (
	"goa-golang/app/controller/billingController"
	"goa-golang/app/controller/userController"
	"goa-golang/app/repository/billingRepository"
	"goa-golang/app/repository/userRepository"
	"goa-golang/app/service/billingService"
	"goa-golang/app/service/userService"
	"goa-golang/internal/logger"
	"goa-golang/internal/storage"

	"github.com/google/wire"
)

func initUserController(db *storage.DbStore, logger logger.Logger) userController.UserControllerInterface {
	wire.Build(userRepository.NewUserRepository, userService.NewUserService, userController.NewUserController)

	return &userController.UserController{}
}

func initBillingController(db *storage.DbStore, uservice userService.UserServiceInterface, logger logger.Logger) billingController.BillingControllerInterface {
	wire.Build(billingRepository.NewBillingRepository, billingService.NewBillingService, billingController.NewBillingController)

	return &billingController.BillingController{}
}

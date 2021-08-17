//go:build wireinject
// +build wireinject

//
// The build tag makes sure the stub is not built in the final build.
//
//lint:file-ignore U1000 Ignore all unused code

package dic

import (
	"goa-golang/app/controller/authController"
	"goa-golang/app/repository/userRepository"
	"goa-golang/app/service/authService"
	"goa-golang/internal/logger"
	"goa-golang/internal/storage"

	"github.com/google/wire"
)

func initUserRepository(db *storage.DbStore) userRepository.UserRepositoryInterface {
	wire.Build(userRepository.NewUserRepository)

	return &userRepository.UserRepository{}
}

// func initUserService(userRepo userRepository.UserRepositoryInterface) userService.UserServiceInterface {
// 	wire.Build(userService.NewUserService)

// 	return &userService.UserService{}
// }

// func initUserController(us userService.UserServiceInterface, logger logger.Logger) userController.UserControllerInterface {
// 	wire.Build(userController.NewUserController)

// 	return &userController.UserController{}
// }

func initAuthService(userRepo userRepository.UserRepositoryInterface, logger logger.Logger) authService.AuthServiceInterface {
	wire.Build(authService.NewAuthService)

	return &authService.AuthService{}
}

func initAuthController(us authService.AuthServiceInterface, logger logger.Logger) authController.AuthControllerInterface {
	wire.Build(authController.NewAuthController)

	return &authController.AuthController{}
}

// func initBillingService(db *storage.DbStore) billingService.BillingServiceInterface {
// 	wire.Build(billingRepository.NewBillingRepository, billingService.NewBillingService)

// 	return &billingService.BillingService{}
// }

// func initBillingController(ub billingService.BillingServiceInterface, us userService.UserServiceInterface, logger logger.Logger) billingController.BillingControllerInterface {
// 	wire.Build(billingController.NewBillingController)

// 	return &billingController.BillingController{}
// }

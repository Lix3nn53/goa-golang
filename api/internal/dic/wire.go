//go:build wireinject
// +build wireinject

//
// The build tag makes sure the stub is not built in the final build.
//
//lint:file-ignore U1000 Ignore all unused code

package dic

import (
	"goa-golang/app/controller/authController"
	"goa-golang/app/controller/characterController"
	"goa-golang/app/controller/mojangController"
	"goa-golang/app/controller/playerController"
	"goa-golang/app/controller/userController"
	"goa-golang/app/repository/characterRepository"
	"goa-golang/app/repository/playerRepository"
	"goa-golang/app/repository/userRepository"
	"goa-golang/app/service/authService"
	"goa-golang/app/service/characterService"
	"goa-golang/app/service/mojangService"
	"goa-golang/app/service/playerService"
	"goa-golang/app/service/userService"
	"goa-golang/internal/logger"
	"goa-golang/internal/storage"

	"github.com/google/wire"
)

// User
func initUserRepository(db *storage.DbStore) userRepository.UserRepositoryInterface {
	wire.Build(userRepository.NewUserRepository)

	return &userRepository.UserRepository{}
}

func initUserService(userRepo userRepository.UserRepositoryInterface) userService.UserServiceInterface {
	wire.Build(userService.NewUserService)

	return &userService.UserService{}
}

func initUserController(us userService.UserServiceInterface, logger logger.Logger) userController.UserControllerInterface {
	wire.Build(userController.NewUserController)

	return &userController.UserController{}
}

// Player
func initPlayerRepository(db *storage.DbStore) playerRepository.PlayerRepositoryInterface {
	wire.Build(playerRepository.NewPlayerRepository)

	return &playerRepository.PlayerRepository{}
}

func initPlayerService(playerRepo playerRepository.PlayerRepositoryInterface) playerService.PlayerServiceInterface {
	wire.Build(playerService.NewPlayerService)

	return &playerService.PlayerService{}
}

func initPlayerController(ps playerService.PlayerServiceInterface, logger logger.Logger) playerController.PlayerControllerInterface {
	wire.Build(playerController.NewPlayerController)

	return &playerController.PlayerController{}
}

// Character
func initCharacterRepository(db *storage.DbStore) characterRepository.CharacterRepositoryInterface {
	wire.Build(characterRepository.NewCharacterRepository)

	return &characterRepository.CharacterRepository{}
}

func initCharacterService(characterRepo characterRepository.CharacterRepositoryInterface) characterService.CharacterServiceInterface {
	wire.Build(characterService.NewCharacterService)

	return &characterService.CharacterService{}
}

func initCharacterController(ps characterService.CharacterServiceInterface, logger logger.Logger) characterController.CharacterControllerInterface {
	wire.Build(characterController.NewCharacterController)

	return &characterController.CharacterController{}
}

// Auth
func initAuthService(playerRepo playerRepository.PlayerRepositoryInterface, userRepo userRepository.UserRepositoryInterface, logger logger.Logger) authService.AuthServiceInterface {
	wire.Build(authService.NewAuthService)

	return &authService.AuthService{}
}

func initAuthController(us authService.AuthServiceInterface, logger logger.Logger) authController.AuthControllerInterface {
	wire.Build(authController.NewAuthController)

	return &authController.AuthController{}
}

// Mojang
func initMojangService(logger logger.Logger) mojangService.MojangServiceInterface {
	wire.Build(mojangService.NewMojangService)

	return &mojangService.MojangService{}
}

func initMojangController(us mojangService.MojangServiceInterface, logger logger.Logger) mojangController.MojangControllerInterface {
	wire.Build(mojangController.NewMojangController)

	return &mojangController.MojangController{}
}

// func initBillingService(db *storage.DbStore) billingService.BillingServiceInterface {
// 	wire.Build(billingRepository.NewBillingRepository, billingService.NewBillingService)

// 	return &billingService.BillingService{}
// }

// func initBillingController(ub billingService.BillingServiceInterface, us userService.UserServiceInterface, logger logger.Logger) billingController.BillingControllerInterface {
// 	wire.Build(billingController.NewBillingController)

// 	return &billingController.BillingController{}
// }

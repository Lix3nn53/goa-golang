// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire

package dic

import (
	"goa-golang/app/controller/userController"
	"goa-golang/app/repository/userRepository"
	"goa-golang/app/service/userService"
	"goa-golang/internal/logger"
	"goa-golang/internal/storage"
)

// Injectors from wire.go:

func InitUserController(db *storage.DbStore, logger2 logger.Logger) userController.UserControllerInterface {
	userRepositoryInterface := userRepository.NewUserRepository(db)
	userServiceInterface := userService.NewUserService(userRepositoryInterface)
	userControllerInterface := userController.NewUserController(userServiceInterface, logger2)
	return userControllerInterface
}

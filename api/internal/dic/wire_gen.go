// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package dic

import (
	"goa-golang/app/controller"
	"goa-golang/app/repository"
	"goa-golang/app/service"
	"goa-golang/internal/logger"
	"goa-golang/internal/storage"
)

// Injectors from wire.go:

func InitUserController(db *storage.DbStore, logger2 logger.Logger) controller.UserControllerInterface {
	userRepositoryInterface := repository.NewUserRepository(db)
	userServiceInterface := service.NewUserService(userRepositoryInterface)
	userControllerInterface := controller.NewUserController(userServiceInterface, logger2)
	return userControllerInterface
}
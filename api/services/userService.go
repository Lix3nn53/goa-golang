package services

import (
	"goa-golang/models"
	"goa-golang/models/scopes"
	"goa-golang/repositories"
	"goa-golang/utils"
)

type IUserService interface {
	GetUsersWithPaginationAndOrder(pagination utils.IPagination, order utils.IOrder) (users []models.User, totalCount int64, err error)
	GetUserByID(id uint) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (service *UserService) GetUserByID(id uint) (user models.User, err error) {
	return service.UserRepository.GetUserByID(id)
}

func (service *UserService) GetUserByEmail(email string) (user models.User, err error) {
	return service.UserRepository.GetUserByEmail(email)
}

func (service *UserService) GetUsersWithPaginationAndOrder(pagination utils.IPagination, order utils.IOrder) (users []models.User, totalCount int64, err error) {
	var userIDs []uint
	userIDs, err = service.UserRepository.GetUserIDs()
	totalCount = int64(len(userIDs))
	users, err = service.UserRepository.GetUsersWithPaginationAndOrder(userIDs, &scopes.GormPagination{Pagination: pagination.Get()}, &scopes.GormOrder{Order: order.Get()})
	return
}

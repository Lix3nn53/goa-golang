package service

import (
	"goa-golang/app/model"
	"goa-golang/app/repository"
)

//UserServiceInterface define the user service interface methods
type UserServiceInterface interface {
	FindByID(id int) (user *model.User, err error)
	RemoveByID(id int) error
	UpdateByID(id int, user model.UpdateUser) error
	Store(user model.CreateUser) (*model.User, error)
}

// billingService handles communication with the user repository
type userService struct {
	userRepo repository.UserRepositoryInterface
}

// NewUserService implements the user service interface.
func NewUserService(userRepo repository.UserRepositoryInterface) *userService {
	return &userService{
		userRepo,
	}
}

// FindByID implements the method to find a user model by primary key
func (s *userService) FindByID(id int) (user *model.User, err error) {
	return s.userRepo.FindByID(id)
}

// FindByID implements the method to remove a user model by primary key
func (s *userService) RemoveByID(id int) error {
	return s.userRepo.RemoveByID(id)
}

// FindByID implements the method to update a user model by primary key
func (s *userService) UpdateByID(id int, user model.UpdateUser) error {
	return s.userRepo.UpdateByID(id, user)
}

// FindByID implements the method to store a new a user model
func (s *userService) Store(user model.CreateUser) (*model.User, error) {
	return s.userRepo.Create(user)
}

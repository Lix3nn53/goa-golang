package userService

import (
	"goa-golang/app/model/userModel"
	"goa-golang/app/repository/userRepository"
)

//UserServiceInterface define the user service interface methods
type UserServiceInterface interface {
	FindByID(uuid string) (user *userModel.User, err error)
	RemoveByID(uuid string) error
	UpdateByID(uuid string, user userModel.UpdateUser) error
	Store(uuid string, user userModel.CreateUser) (*userModel.User, error)
}

// billingService handles communication with the user repository
type UserService struct {
	userRepo userRepository.UserRepositoryInterface
}

// NewUserService implements the user service interface.
func NewUserService(userRepo userRepository.UserRepositoryInterface) UserServiceInterface {
	return &UserService{
		userRepo,
	}
}

// FindByID implements the method to find a user model by primary key
func (s *UserService) FindByID(uuid string) (user *userModel.User, err error) {
	return s.userRepo.FindByID(uuid)
}

// FindByID implements the method to remove a user model by primary key
func (s *UserService) RemoveByID(uuid string) error {
	return s.userRepo.RemoveByID(uuid)
}

// FindByID implements the method to update a user model by primary key
func (s *UserService) UpdateByID(uuid string, user userModel.UpdateUser) error {
	return s.userRepo.UpdateByID(uuid, user)
}

// FindByID implements the method to store a new a user model
func (s *UserService) Store(uuid string, user userModel.CreateUser) (*userModel.User, error) {
	return s.userRepo.Create(uuid, user)
}

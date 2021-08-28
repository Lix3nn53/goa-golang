package userService

import (
	"goa-golang/app/model/userModel"
	"goa-golang/app/repository/userRepository"
)

//UserServiceInterface define the user service interface methods
type UserServiceInterface interface {
	FindByID(id string, field string) (user *userModel.User, err error)
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
func (s *UserService) FindByID(id string, field string) (user *userModel.User, err error) {
	return s.userRepo.FindByID(id, field)
}

package authService

import (
	"goa-golang/app/model/userModel"
	"goa-golang/app/repository/authRepository"
)

//UserServiceInterface define the user service interface methods
type AuthServiceInterface interface {
	GoogleOauth2() (user *userModel.User, err error)
}

// billingService handles communication with the user repository
type AuthService struct {
	authRepo authRepository.AuthRepositoryInterface
}

// NewUserService implements the user service interface.
func NewAuthService(authRepo authRepository.AuthRepositoryInterface) AuthServiceInterface {
	return &AuthService{
		authRepo,
	}
}

// FindByID implements the method to find a user model by primary key
func (s *AuthService) GoogleOauth2() (user *userModel.User, err error) {
	return s.authRepo.GoogleOauth2()
}

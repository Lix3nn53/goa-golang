package policies

import (
	"goa-golang/models"
)

type IUserPolicy interface {
	Index(auth models.User) bool
	Show(auth models.User, user models.User) bool
	Update(auth models.User, user models.User) bool
	Delete(auth models.User, user models.User) bool
}

type UserPolicy struct{}

func (UserPolicy) Index(auth models.User) bool {
	return auth.Admin
}

func (UserPolicy) Show(auth models.User, user models.User) bool {
	return auth.ID == 1 && user.Verified
}

func (UserPolicy) Update(auth models.User, user models.User) bool {
	return auth.ID == user.ID && user.Verified
}

func (UserPolicy) Delete(auth models.User, user models.User) bool {
	return auth.ID == user.ID && user.Verified
}

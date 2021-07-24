package requests

import (
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type LoginRequest struct {
	Email    string `json:"email" form:"email" xml:"email" binding:"required,email"`
	Password string `json:"password" form:"password" xml:"password" binding:"required"`
}

func (r LoginRequest) RegisterValidation() {
	var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
		date, ok := fl.Field().Interface().(time.Time)
		if ok {
			today := time.Now()
			if today.After(date) {
				return false
			}
		}
		return true
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}
}

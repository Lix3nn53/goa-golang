package userModel

// User represents user resources.
type User struct {
	UUID       string `json:"uuid"`
	Email      string `json:"email"`
	McUsername string `json:"mc_username"`
	Credits    int    `json:"credits"`
}

// UpdateUser represents user resources.
type UpdateUser struct {
	Email      string `json:"email" validate:"required"`
	McUsername string `json:"mc_username" validate:"required"`
	Credits    int    `json:"credits" validate:"required"`
}

// CreateUser represents user resources.
type CreateUser struct {
	Email      string `json:"email" validate:"required"`
	McUsername string `json:"mc_username" validate:"required"`
	Credits    int    `json:"credits" validate:"required"`
}

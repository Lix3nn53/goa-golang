package userModel

// User represents user resources.
type User struct {
	UUID    string `json:"uuid" db:"uuid"`
	Email   string `json:"email" db:"email"`
	Credits int    `json:"credits" db:"credits"`
}

// UpdateUser represents user resources.
type UpdateUser struct {
	Email   string `json:"email" validate:"required"`
	Credits int    `json:"credits" validate:"required"`
}

// CreateUser represents user resources.
type CreateUserMicrosoft struct {
	UUID string `json:"uuid" validate:"required"`
}

// CreateUser represents user resources.
type CreateUserGoogle struct {
	GoogleId string `json:"google_id" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

// UpdateUser represents user resources.
type Sessions struct {
	Sessions string `json:"sessions" validate:"required"`
}

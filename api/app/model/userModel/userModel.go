package userModel

// User represents user resources.
type User struct {
	UUID       string `json:"uuid" db:"uuid"`
	Email      string `json:"email" db:"email"`
	McUsername string `json:"mc_username" db:"mc_username"`
	Credits    int    `json:"credits" db:"credits"`
	Sessions   string `json:"sessions" db:"sessions"`
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

// UpdateUser represents user resources.
type UserSessions struct {
	Sessions string `json:"sessions" validate:"required"`
}

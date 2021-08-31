package userModel

import "database/sql"

// User represents user resources.
type User struct {
	UUID      string `json:"uuid" db:"uuid"`
	GoogleId  string `json:"google_id" db:"google_id"`
	DiscordId string `json:"discord_id" db:"discord_id"`
	Email     string `json:"email" db:"email"`
	Credits   int    `json:"credits" db:"credits"`
}

type UserScan struct {
	UUID      sql.NullString `json:"uuid" db:"uuid"`
	GoogleId  sql.NullString `json:"google_id" db:"google_id"`
	DiscordId sql.NullString `json:"discord_id" db:"discord_id"`
	Email     sql.NullString `json:"email" db:"email"`
	Credits   sql.NullInt16  `json:"credits" db:"credits"`
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

// CreateUser represents user resources.
type CreateUserDiscord struct {
	DiscordId string `json:"discord_id" validate:"required"`
}

// UpdateUser represents user resources.
type Sessions struct {
	Sessions string `json:"sessions" validate:"required"`
}

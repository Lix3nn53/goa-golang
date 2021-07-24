package model

// User represents user resources.
type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Cif        string `json:"cif"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

// UpdateUser represents user resources.
type UpdateUser struct {
	Name       string `json:"name" validate:"required"`
	Cif        string `json:"cif" validate:"required"`
	Country    string `json:"country" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
}

// CreateUser represents user resources.
type CreateUser struct {
	Name       string `json:"name" validate:"required"`
	Cif        string `json:"cif" validate:"required"`
	Country    string `json:"country" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
}

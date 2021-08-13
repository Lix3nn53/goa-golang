package authRepository

import (
	"database/sql"
	appError "goa-golang/app/error"
	"goa-golang/app/model/userModel"
	"goa-golang/internal/storage"
)

// billingRepository handles communication with the user store
type authRepository struct {
	db *storage.DbStore
}

//UserRepositoryInterface define the user repository interface methods
type AuthRepositoryInterface interface {
	GoogleOauth2(id int) (user *userModel.User, err error)
}

// NewUserRepository implements the user repository interface.
func NewAuthRepository(db *storage.DbStore) AuthRepositoryInterface {
	return &authRepository{
		db,
	}
}

// FindByID implements the method to find a user from the store
func (r *authRepository) GoogleOauth2(id int) (user *userModel.User, err error) {
	user = &userModel.User{}

	var query = "SELECT id, cif, name, postal_code, country FROM users WHERE id = $1"
	row := r.db.QueryRow(query, id)

	if err := row.Scan(&user.ID, &user.Cif, &user.Name, &user.PostalCode, &user.Country); err != nil {
		if err == sql.ErrNoRows {
			return nil, appError.ErrNotFound
		}

		return nil, err
	}

	return user, nil
}

package userRepository

import (
	"database/sql"
	appError "goa-golang/app/error"
	"goa-golang/app/model/userModel"
	"goa-golang/internal/storage"
)

// billingRepository handles communication with the user store
type userRepository struct {
	db *storage.DbStore
}

//UserRepositoryInterface define the user repository interface methods
type UserRepositoryInterface interface {
	FindByID(uuid string) (user *userModel.User, err error)
	RemoveByID(uuid string) error
	UpdateByID(uuid string, user userModel.UpdateUser) error
	Create(uuid string, create userModel.CreateUser) (user *userModel.User, err error)
}

// NewUserRepository implements the user repository interface.
func NewUserRepository(db *storage.DbStore) UserRepositoryInterface {
	return &userRepository{
		db,
	}
}

// FindByID implements the method to find a user from the store
func (r *userRepository) FindByID(uuid string) (user *userModel.User, err error) {
	user = &userModel.User{}

	var query = "SELECT id, cif, name, postal_code, country FROM users WHERE id = $1"
	row := r.db.QueryRow(query, uuid)

	if err := row.Scan(&user.UUID, &user.Email, &user.McUsername, &user.Credits); err != nil {
		if err == sql.ErrNoRows {
			return nil, appError.ErrNotFound
		}

		return nil, err
	}

	return user, nil
}

// RemoveByID implements the method to remove a user from the store
func (r *userRepository) RemoveByID(uuid string) error {

	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1;`, uuid)
	return err
}

// UpdateByID implements the method to update a user into the store
func (r *userRepository) UpdateByID(uuid string, user userModel.UpdateUser) error {
	result, err := r.db.Exec("UPDATE users SET email = $1, mc_username = $2, credits = $3 where id = $4", user.Email, user.McUsername, user.Credits, uuid)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows != 1 {
		return appError.ErrNotFound
	}

	return nil
}

// Create implements the method to persist a new user
func (r *userRepository) Create(uuid string, UserSignUp userModel.CreateUser) (user *userModel.User, err error) {
	createUserQuery := `INSERT INTO goa_web_player (uuid, email, mc_username, credits) 
		VALUES ($1, $2, $3, $4)`

	stmt, err := r.db.Prepare(createUserQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var userModel userModel.User
	err = stmt.QueryRow(uuid, UserSignUp.Email, UserSignUp.McUsername, UserSignUp.Credits).Scan(&userModel)
	if err != nil {
		return nil, err
	}

	return &userModel, nil
}

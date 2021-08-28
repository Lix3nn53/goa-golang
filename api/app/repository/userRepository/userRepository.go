package userRepository

import (
	"database/sql"
	"errors"
	appError "goa-golang/app/error"
	"goa-golang/app/model/userModel"
	"goa-golang/internal/storage"
	"strings"
)

// billingRepository handles communication with the user store
type UserRepository struct {
	db *storage.DbStore
}

//UserRepositoryInterface define the user repository interface methods
type UserRepositoryInterface interface {
	FindByID(id string, field string) (user *userModel.User, err error)
	RemoveByID(uuid string) error
	UpdateByID(uuid string, user userModel.UpdateUser) error
	CreateWithMicrosoft(create userModel.CreateUserMicrosoft) (user *userModel.User, err error)
	CreateWithGoogle(create userModel.CreateUserGoogle) (user *userModel.User, err error)
	GetSessions(uuid string) (sessions string, err error)
	AddSession(uuid string, refreshToken string) error
	RemoveSession(uuid string, refreshToken string) error
}

// NewUserRepository implements the user repository interface.
func NewUserRepository(db *storage.DbStore) UserRepositoryInterface {
	return &UserRepository{
		db,
	}
}

// FindByID implements the method to find a user from the store
func (r *UserRepository) FindByID(id string, field string) (user *userModel.User, err error) {
	user = &userModel.User{}

	var query = "SELECT uuid, email, credits FROM goa_player_web WHERE ? = ?"
	row := r.db.QueryRow(query, field, id)

	if err := row.Scan(&user.UUID, &user.Email, &user.Credits); err != nil {
		return nil, err
	}

	return user, nil
}

// FindByGoogle implements the method to find a user from the store
func (r *UserRepository) FindByGoogle(googleId string) (user *userModel.User, err error) {
	user = &userModel.User{}

	var query = "SELECT uuid, email, credits FROM goa_player_web WHERE google_id = ?"
	row := r.db.QueryRow(query, googleId)

	if err := row.Scan(&user.UUID, &user.Email, &user.Credits); err != nil {
		return nil, err
	}

	return user, nil
}

// RemoveByID implements the method to remove a user from the store
func (r *UserRepository) RemoveByID(uuid string) error {

	_, err := r.db.Exec(`DELETE FROM goa_player_web WHERE uuid = ?;`, uuid)
	return err
}

// UpdateByID implements the method to update a user into the store
func (r *UserRepository) UpdateByID(uuid string, user userModel.UpdateUser) error {
	result, err := r.db.Exec("UPDATE goa_player_web SET email = ?, credits = ? where uuid = ?", user.Email, user.Credits, uuid)
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
func (r *UserRepository) CreateWithMicrosoft(UserSignUp userModel.CreateUserMicrosoft) (user *userModel.User, err error) {
	createUserQuery := `INSERT INTO goa_player_web (uuid) VALUES (?)`

	stmt, err := r.db.Prepare(createUserQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(UserSignUp.UUID)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	n := int(rows) // truncated on machines with 32-bit ints
	if n == 0 {
		return nil, appError.ErrNotFound
	}

	return &userModel.User{
		UUID:    UserSignUp.UUID,
		Credits: 0,
	}, nil
}

// Create implements the method to persist a new user
func (r *UserRepository) CreateWithGoogle(UserSignUp userModel.CreateUserGoogle) (user *userModel.User, err error) {
	createUserQuery := `INSERT INTO goa_player_web (google_id, email) VALUES (?, ?)`

	stmt, err := r.db.Prepare(createUserQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(UserSignUp.GoogleId, UserSignUp.Email)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	n := int(rows) // truncated on machines with 32-bit ints
	if n == 0 {
		return nil, appError.ErrNotFound
	}

	return &userModel.User{
		Email:   UserSignUp.Email,
		Credits: 0,
	}, nil
}

// FindByID implements the method to find a user from the store
func (r *UserRepository) GetSessions(uuid string) (sessions string, err error) {
	var query = "SELECT sessions FROM goa_player_web WHERE uuid = ?"
	row := r.db.QueryRow(query, uuid)

	var sessionsScan sql.NullString
	if err := row.Scan(&sessionsScan); err != nil {
		return "", err
	}
	if !sessionsScan.Valid {
		return "", errors.New("sql string is not valid")
	}

	return sessionsScan.String, nil
}

func (r *UserRepository) AddSession(uuid string, refreshToken string) error {
	sessions, err := r.GetSessions(uuid)
	if err != nil {
		return err
	}

	sessions = sessions + "/" + refreshToken

	result, err := r.db.Exec("UPDATE goa_player_web SET sessions = ? where uuid = ?", sessions, uuid)
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

func (r *UserRepository) RemoveSession(uuid string, refreshToken string) error {
	sessions, err := r.GetSessions(uuid)
	if err != nil {
		return err
	}

	sessions = strings.Replace(sessions, "/"+refreshToken, "", -1)

	result, err := r.db.Exec("UPDATE goa_player_web SET sessions = ? where uuid = ?", sessions, uuid)
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

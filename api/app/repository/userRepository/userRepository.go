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
	UpdateByID(id string, field string, user userModel.UpdateUser) error
	CreateWithMicrosoft(create userModel.CreateUserMicrosoft) (user *userModel.User, err error)
	CreateWithGoogle(create userModel.CreateUserGoogle) (user *userModel.User, err error)
	CreateWithDiscord(create userModel.CreateUserDiscord) (user *userModel.User, err error)
	GetSessions(id string, field string) (sessions sql.NullString)
	AddSession(id string, field string, refreshToken string) error
	RemoveSession(id string, field string, refreshToken string) error
}

// NewUserRepository implements the user repository interface.
func NewUserRepository(db *storage.DbStore) UserRepositoryInterface {
	return &UserRepository{
		db,
	}
}

// FindByID implements the method to find a user from the store
func (r *UserRepository) FindByID(id string, field string) (user *userModel.User, err error) {
	var query = "SELECT uuid, google_id, discord_id, email, credits FROM goa_player_web WHERE " + field + " = ?"
	row := r.db.QueryRow(query, id)

	scan := &userModel.UserScan{}

	if err := row.Scan(&scan.UUID, &scan.GoogleId, &scan.DiscordId, &scan.Email, &scan.Credits); err != nil {
		return nil, err
	}

	user = &userModel.User{
		UUID:      scan.UUID.String,
		GoogleId:  scan.GoogleId.String,
		DiscordId: scan.DiscordId.String,
		Email:     scan.Email.String,
		Credits:   int(scan.Credits.Int16),
	}

	return user, nil
}

// RemoveByID implements the method to remove a user from the store
func (r *UserRepository) RemoveByID(uuid string) error {

	_, err := r.db.Exec(`DELETE FROM goa_player_web WHERE uuid = ?;`, uuid)
	return err
}

// UpdateByID implements the method to update a user into the store
func (r *UserRepository) UpdateByID(id string, field string, user userModel.UpdateUser) error {
	result, err := r.db.Exec("UPDATE goa_player_web SET email = ?, credits = ? where "+field+" = ?", user.Email, user.Credits, id)
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
		GoogleId: UserSignUp.GoogleId,
		Email:    UserSignUp.Email,
		Credits:  0,
	}, nil
}

// Create implements the method to persist a new user
func (r *UserRepository) CreateWithDiscord(UserSignUp userModel.CreateUserDiscord) (user *userModel.User, err error) {
	createUserQuery := `INSERT INTO goa_player_web (discord_id) VALUES (?)`

	stmt, err := r.db.Prepare(createUserQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(UserSignUp.DiscordId)
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
		GoogleId: UserSignUp.DiscordId,
		Credits:  0,
	}, nil
}

// FindByID implements the method to find a user from the store
func (r *UserRepository) GetSessions(id string, field string) (sessions sql.NullString) {
	var query = "SELECT sessions FROM goa_player_web WHERE " + field + " = ?"
	row := r.db.QueryRow(query, id)

	var scan sql.NullString

	if err := row.Scan(&scan); err != nil {
		return sql.NullString{Valid: false}
	}

	return scan
}

func (r *UserRepository) AddSession(id string, field string, refreshToken string) error {
	sessionsStr := refreshToken

	sessions := r.GetSessions(id, field)
	if sessions.Valid {
		sessionsStr = sessions.String + "/" + refreshToken
	}

	result, err := r.db.Exec("UPDATE goa_player_web SET sessions = ? where "+field+" = ?", sessionsStr, id)
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

func (r *UserRepository) RemoveSession(id string, field string, refreshToken string) error {
	sessions := r.GetSessions(id, field)
	if !sessions.Valid {
		return errors.New("sessions is not valid")
	}

	sessionsStr := strings.Replace(sessions.String, "/"+refreshToken, "", -1)

	result, err := r.db.Exec("UPDATE goa_player_web SET sessions = ? where "+field+" = ?", sessionsStr, id)
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

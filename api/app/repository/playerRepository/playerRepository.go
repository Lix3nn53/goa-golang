package playerRepository

import (
	"database/sql"
	appError "goa-golang/app/error"
	"goa-golang/app/model/playerModel"
	"goa-golang/internal/storage"
)

// billingRepository handles communication with the user store
type PlayerRepository struct {
	db *storage.DbStore
}

//PlayerRepositoryInterface define the user repository interface methods
type PlayerRepositoryInterface interface {
	FindByID(uuid string) (user *playerModel.Player, err error)
	RemoveByID(uuid string) error
	CreateUUID(createPlayer playerModel.CreatePlayer) (err error)
}

// NewPlayerRepository implements the user repository interface.
func NewPlayerRepository(db *storage.DbStore) PlayerRepositoryInterface {
	return &PlayerRepository{
		db,
	}
}

// FindByID implements the method to find a user from the store
func (r *PlayerRepository) FindByID(uuid string) (player *playerModel.Player, err error) {
	var query = "SELECT daily_last_date, staff_rank, premium_rank, premium_rank_date FROM goa_player WHERE uuid = ?"
	row := r.db.QueryRow(query, uuid)

	var read1 sql.NullString
	var read2 sql.NullString
	var read3 sql.NullString
	var read4 sql.NullString
	if err := row.Scan(&read1, &read2, &read3, &read4); err != nil {
		return nil, err
	}

	var dailyLastDate string
	if read1.Valid {
		dailyLastDate = read1.String
	}
	var staffRank string
	if read2.Valid {
		staffRank = read2.String
	}
	var premiumRank string
	if read3.Valid {
		premiumRank = read3.String
	}
	var premiumRankDate string
	if read4.Valid {
		premiumRankDate = read4.String
	}

	player = &playerModel.Player{
		UUID:            uuid,
		DailyLastDate:   dailyLastDate,
		StaffRank:       staffRank,
		PremiumRank:     premiumRank,
		PremiumRankDate: premiumRankDate,
	}

	return player, nil
}

// RemoveByID implements the method to remove a user from the store
func (r *PlayerRepository) RemoveByID(uuid string) error {

	_, err := r.db.Exec(`DELETE FROM goa_player WHERE uuid = ?;`, uuid)
	return err
}

// Create implements the method to persist a new user
func (r *PlayerRepository) CreateUUID(createPlayer playerModel.CreatePlayer) (err error) {
	createUserQuery := `INSERT INTO goa_player (uuid) 
		VALUES (?)`

	stmt, err := r.db.Prepare(createUserQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(createPlayer.UUID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	n := int(rows) // truncated on machines with 32-bit ints
	if n == 0 {
		return appError.ErrNotFound
	}

	return nil
}

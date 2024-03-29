package playerRepository

import (
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

	scan := &playerModel.PlayerScan{}

	if err := row.Scan(&scan.DailyLastDate, &scan.StaffRank, &scan.PremiumRank, &scan.PremiumRankDate); err != nil {
		return nil, err
	}

	player = &playerModel.Player{
		UUID:            uuid,
		DailyLastDate:   scan.DailyLastDate.String,
		StaffRank:       scan.StaffRank.String,
		PremiumRank:     scan.PremiumRank.String,
		PremiumRankDate: scan.PremiumRankDate.String,
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
	createUserQuery := `INSERT INTO goa_player (uuid)	VALUES (?)`

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

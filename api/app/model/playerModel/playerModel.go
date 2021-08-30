package playerModel

import "database/sql"

// User represents user resources.
type Player struct {
	UUID            string `json:"uuid" db:"uuid"`
	DailyLastDate   string `json:"daily_last_date" db:"daily_last_date"`
	StaffRank       string `json:"staff_rank" db:"staff_rank"`
	PremiumRank     string `json:"premium_rank" db:"premium_rank"`
	PremiumRankDate string `json:"premium_rank_date" db:"premium_rank_date"`
}

type PlayerScan struct {
	UUID            sql.NullString `json:"uuid" db:"uuid"`
	DailyLastDate   sql.NullString `json:"daily_last_date" db:"daily_last_date"`
	StaffRank       sql.NullString `json:"staff_rank" db:"staff_rank"`
	PremiumRank     sql.NullString `json:"premium_rank" db:"premium_rank"`
	PremiumRankDate sql.NullString `json:"premium_rank_date" db:"premium_rank_date"`
}

// CreateUser represents user resources.
type CreatePlayer struct {
	UUID string `json:"uuid" db:"uuid" validate:"required"`
}

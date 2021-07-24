package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

// DbStore ...
type DbStore struct {
	*sqlx.DB
}

// InitializeDB Opening a storage and save the reference to `Database` struct.
func InitializeDB() *DbStore {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable TimeZone=Europe/Paris password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		return nil
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	retryCount := 30
	for {
		err := db.Ping()
		if err != nil {
			if retryCount == 0 {
				log.Fatalf("Not able to establish connection to database")
			}
			log.Printf(fmt.Sprintf("Could not connect to database. Wait 2 seconds. %d retries left...", retryCount))
			retryCount--
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	if err = db.Ping(); err != nil {
		return nil
	}

	return &DbStore{
		db,
	}
}

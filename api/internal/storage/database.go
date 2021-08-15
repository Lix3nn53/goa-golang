package storage

import (
	"fmt"
	"goa-golang/internal/logger"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
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
func InitializeDB(logger logger.Logger) *DbStore {
	db, err := sqlx.Connect("mysql", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		logger.Fatalf(err.Error())
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
				logger.Fatalf("Not able to establish connection to database")
			}
			logger.Infof(fmt.Sprintf("Could not connect to database. Wait 2 seconds. %d retries left...", retryCount))
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

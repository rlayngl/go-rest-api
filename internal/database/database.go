package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)

	return db, nil
}

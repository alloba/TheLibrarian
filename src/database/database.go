package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(connectionString string) *sql.DB {
	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		panic(fmt.Errorf("unable to connect to database: %w", err))
	}

	return db
}

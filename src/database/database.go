package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(connectionString string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("unable to connect to database: %w", err))
	}

	return db
}

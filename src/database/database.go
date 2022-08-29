package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func Connect(connectionString string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("unable to connect to database: %w", err))
	}

	if err = db.Exec("PRAGMA foreign_keys = ON", nil).Error; err != nil {
		log.Fatalf("could not enable foreign key support. failing outright.")
	}

	return db
}

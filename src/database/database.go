package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func Connect(connectionString string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Panicf("Connect: unable to connect to database; %v", err.Error())
	}

	if err = db.Exec("PRAGMA foreign_keys = ON", nil).Error; err != nil {
		log.Panicf("Connect: enable to enable foreign key support; %v", err.Error())
	}
	return db
}

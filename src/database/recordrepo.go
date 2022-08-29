package database

import (
	"fmt"
	"gorm.io/gorm"
)

func NewRecordRepo(db *gorm.DB) *RecordRepo {
	return &RecordRepo{
		db:         db,
		FindAll:    getAllRecords(db),
		SaveOne:    saveSingleRecord(db),
		FindByHash: findSingleRecordByHash(db),
	}
}

type RecordRepo struct {
	db         *gorm.DB
	FindAll    findAllRecords
	SaveOne    saveRecord
	FindByHash findRecordByHash
}

type findAllRecords func() ([]Record, error)
type saveRecord func(record *Record) error
type findRecordByHash func(hash string) (*Record, error)

func getAllRecords(db *gorm.DB) func() ([]Record, error) {
	return func() ([]Record, error) {
		var records []Record
		result := db.Find(&records)
		if result.Error != nil {
			return nil, fmt.Errorf("could not get all records from database - %v", result.Error.Error())
		}
		return records, nil
	}
}

func findSingleRecordByHash(db *gorm.DB) func(hash string) (*Record, error) {
	return func(hash string) (*Record, error) {
		var record Record
		result := db.Where("hash = ?", hash).First(&record)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to find record - %v", result.Error.Error())
		}
		return &record, nil
	}
}

func saveSingleRecord(db *gorm.DB) func(record *Record) error {
	return func(record *Record) error {
		result := db.Create(record)
		if result.Error != nil {
			return fmt.Errorf("could not save record to database - %v", result.Error.Error())
		}
		return nil
	}
}

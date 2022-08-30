package database

import (
	"fmt"
	"gorm.io/gorm"
)

func NewRecordRepo(db *gorm.DB) *RecordRepo {
	return &RecordRepo{
		db: db,
	}
}

type RecordRepo struct {
	db *gorm.DB
}

func (repo *RecordRepo) FindAll() ([]Record, error) {
	var records []Record
	result := repo.db.Find(&records)
	if result.Error != nil {
		return nil, fmt.Errorf("could not get all records from database - %v", result.Error.Error())
	}
	return records, nil
}

func (repo *RecordRepo) FindByHash(hash string) (*Record, error) {
	var record Record
	result := repo.db.Where("id = ?", hash).First(&record)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find record - %v", result.Error.Error())
	}
	return &record, nil
}

func (repo *RecordRepo) SaveOne(record *Record) error {
	result := repo.db.Create(record)
	if result.Error != nil {
		return fmt.Errorf("could not save record to database - %v", result.Error.Error())
	}
	return nil
}

func (repo *RecordRepo) Exists(hash string) (bool, error) {
	var exists bool
	err := repo.db.Model(Record{}).
		Select("count(*) > 0").
		Where("id = ?", hash).
		Find(&exists).Error
	if err != nil {
		return false, fmt.Errorf("failed to search for record exsts %v - %v", hash, err.Error())
	}
	return exists, nil
}

package repository

import (
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
)

type RecordRepo struct {
	db *gorm.DB
}

func NewRecordRepo(db *gorm.DB) *RecordRepo {
	return &RecordRepo{db: db}
}

func (repo RecordRepo) FindOne(recordId string) (*database.Record, error) {
	res := database.Record{}
	err := repo.db.Where("id = ?", recordId).First(&res).Error
	if err != nil {
		return nil, logTrace(err)
	}
	return &res, nil
}

func (repo RecordRepo) Exists(recordId string) (bool, error) {
	var exists bool
	err := repo.db.Model(&database.Record{}).
		Select("count(*) > 0").
		Where("id = ?", recordId).
		Find(&exists).Error

	if err != nil {
		return false, logTrace(err)
	}
	return exists, nil
}

func (repo RecordRepo) CreateOne(record *database.Record) error {
	if err := repo.db.Create(record).Error; err != nil {
		return logTrace(err)
	}
	return nil
}

func (repo RecordRepo) UpsertAll(records *[]database.Record) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		for _, rec := range *records {
			err := tx.Save(rec).Error
			if err != nil {
				return logTrace(err)
			}
		}
		return nil
	})
	if err != nil {
		return logTrace(err)
	}
	return nil
}

func (repo RecordRepo) DeleteOne(recordId string) error {
	exist, err := repo.Exists(recordId)
	if err != nil {
		return logTrace(err)
	}

	if !exist {
		return nil
	}

	err = repo.db.Where("id = ?", recordId).Delete(&database.Record{}).Error
	if err != nil {
		return logTrace(err)
	}
	return nil
}

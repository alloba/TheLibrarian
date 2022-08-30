package repository

import (
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
)

type PageRepo struct {
	db *gorm.DB
}

func NewPageRepo(db *gorm.DB) *PageRepo {
	return &PageRepo{db: db}
}

func (repo PageRepo) Exists(PageId string) (bool, error) {
	var exists bool
	err := repo.db.Model(&database.Page{}).
		Select("count(*) > 0").
		Where("id = ?", PageId).
		Find(&exists).Error

	if err != nil {
		return false, logTrace(err)
	}
	return exists, nil
}

func (repo PageRepo) SaveOne(Page *database.Page) error {
	if err := repo.db.Create(Page).Error; err != nil {
		return logTrace(err)
	}
	return nil
}

func (repo PageRepo) DeleteOne(PageId string) error {
	exist, err := repo.Exists(PageId)
	if err != nil {
		return logTrace(err)
	}

	if !exist {
		return nil
	}

	err = repo.db.Where("id = ?", PageId).Delete(&database.Page{}).Error
	if err != nil {
		return logTrace(err)
	}
	return nil
}

package repository

import (
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/logging"
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
		return false, logging.LogTrace(err)
	}
	return exists, nil
}

func (repo PageRepo) ExistsByRecordAndEdition(recordId string, editionId string) (bool, error) {
	var exists bool
	err := repo.db.Model(&database.Page{}).
		Select("count(*) > 0").
		Where("record_id = ? and edition_id = ?", recordId, editionId).
		Find(&exists).Error

	if err != nil {
		return false, logging.LogTrace(err)
	}
	return exists, nil
}

func (repo PageRepo) FindOneByRecordAndEdition(recordId string, editionId string) (*database.Page, error) {
	res := database.Page{}
	err := repo.db.Where("record_id = ? and edition_id = ?", recordId, editionId).First(&res).Error
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	return &res, nil
}

func (repo PageRepo) FindAllByEditionId(editionId string) (*[]database.Page, error) {
	res := make([]database.Page, 0)
	err := repo.db.Where("edition_id = ?", editionId).Find(&res).Error
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	return &res, nil
}

func (repo PageRepo) CreateOne(Page *database.Page) error {
	if err := repo.db.Create(Page).Error; err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

func (repo PageRepo) UpsertAll(pages *[]database.Page) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		for _, page := range *pages {
			err := tx.Save(page).Error
			if err != nil {
				return logging.LogTrace(err)
			}
		}
		return nil
	})
	if err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

func (repo PageRepo) DeleteOne(PageId string) error {
	exist, err := repo.Exists(PageId)
	if err != nil {
		return logging.LogTrace(err)
	}

	if !exist {
		return nil
	}

	err = repo.db.Where("id = ?", PageId).Delete(&database.Page{}).Error
	if err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

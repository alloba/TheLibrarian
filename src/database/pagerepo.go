package database

import (
	"fmt"
	"gorm.io/gorm"
)

type PageRepo struct {
	db *gorm.DB
}

func NewPageRepo(db *gorm.DB) *PageRepo {
	return &PageRepo{
		db: db,
	}
}

func (repo PageRepo) getByUuid(uuid string) (*Page, error) {
	out := Page{}
	err := repo.db.Where("uuid = ?", uuid).Find(&out).Error
	if err != nil {
		return &out, fmt.Errorf("could not get page from database %v", err.Error())
	}
	return &out, nil
}
func (repo PageRepo) getByEditionUuidAndRecordHash(editionUuid string, recordHash string) (*Page, error) {
	out := Page{}
	err := repo.db.Where("editionUuid = ? and recordHash = ?", editionUuid, recordHash).Find(&out).Error
	if err != nil {
		return &out, fmt.Errorf("could not get page from database %v", err.Error())
	}
	return &out, nil
}

func (repo PageRepo) ExistsById(pageUuid string) (bool, error) {
	var exists bool
	err := repo.db.Model(&Page{}).
		Select("count(*) > 1").
		Where("uuid = ?", pageUuid).
		Find(exists).
		Error
	if err != nil {
		return false, fmt.Errorf("could not check if page exists - %v", err.Error())
	}
	return exists, nil
}

func (repo PageRepo) ExistsByEditionUuidAndRecordHash(editionUuid string, recordHash string) (bool, error) {
	var exists bool
	err := repo.db.Model(&Page{}).
		Select("count(*) > 1").
		Where("recordHash = ? and editionUuid = ?", recordHash, editionUuid).
		Find(exists).
		Error
	if err != nil {
		return false, fmt.Errorf("could not check if page exists - %v", err.Error())
	}
	return exists, nil
}

func (repo PageRepo) SaveOne(page *Page) error {
	err := repo.db.Create(page).Error
	if err != nil {
		return fmt.Errorf("could not save page to database - %v", err.Error())
	}
	return nil
}

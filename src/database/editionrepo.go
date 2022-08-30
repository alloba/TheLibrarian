package database

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EditionRepo struct {
	db *gorm.DB
}

func NewEditionRepo(db *gorm.DB) *EditionRepo {
	return &EditionRepo{db: db}
}

func (repo EditionRepo) SaveOne(edition *Edition) error {
	err := repo.db.Preload(clause.Associations).Create(&edition).Error
	if err != nil {
		return fmt.Errorf("could not save edition to database - %v", err.Error())
	}
	return nil
}

func (repo EditionRepo) FindAll() (*[]Edition, error) {
	var editions []Edition
	err := repo.db.Find(&editions).Error
	if err != nil {
		return nil, fmt.Errorf("could not load all editions - %v", err.Error())
	}
	return &editions, nil
}

func (repo EditionRepo) FindOne(uuid string) (*Edition, error) {
	var edition = &Edition{}
	err := repo.db.Preload(clause.Associations).Where("id = ?", uuid).First(&edition).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find edition in db %v - %v", uuid, err.Error())
	}
	return edition, nil
}

func (repo EditionRepo) Exists(uuid string) (bool, error) {
	var exists bool
	err := repo.db.Model(Edition{}).
		Select("count(*) > 0").
		Where("id = ?", uuid).
		Find(&exists).Error
	if err != nil {
		return false, fmt.Errorf("failed to search for edition exsts %v - %v", uuid, err.Error())
	}
	return exists, nil
}

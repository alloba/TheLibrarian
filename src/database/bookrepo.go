package database

import (
	"fmt"
	"gorm.io/gorm"
)

type BookRepo struct {
	db *gorm.DB
}

func NewBookRepo(db *gorm.DB) *BookRepo {
	return &BookRepo{db: db}
}

func (repo BookRepo) FindOne(uuid string) (*Book, error) {
	var res = Book{}
	err := repo.db.Where("uuid = ?", uuid).First(&res).Error
	if err != nil {
		return nil, fmt.Errorf("unable to find book in database: %v - %v", uuid, err.Error())
	}

	return &res, nil
}

func (repo BookRepo) SaveOne(book *Book) error {
	err := repo.db.Create(book).Error
	if err != nil {
		return fmt.Errorf("could not save book - %v", err.Error())
	}
	return nil
}

func (repo BookRepo) Exists(uuid string) (bool, error) {
	var exists bool
	err := repo.db.Model(Book{}).
		Select("count(*) > 0").
		Where("uuid = ?", uuid).
		Find(&exists).Error
	if err != nil {
		return false, fmt.Errorf("failed to search for book exsts %v - %v", uuid, err.Error())
	}
	return exists, nil
}

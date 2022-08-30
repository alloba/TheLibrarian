package repository

import (
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
)

type BookRepo struct {
	db *gorm.DB
}

func NewBookRepo(db *gorm.DB) *BookRepo {
	return &BookRepo{db: db}
}

func (repo BookRepo) Exists(bookId string) (bool, error) {
	var exists bool
	err := repo.db.Model(&database.Book{}).
		Select("count(*) > 0").
		Where("id = ?", bookId).
		Find(&exists).Error

	if err != nil {
		return false, logTrace(err)
	}
	return exists, nil
}

func (repo BookRepo) SaveOne(book *database.Book) error {
	if err := repo.db.Create(book).Error; err != nil {
		return logTrace(err)
	}
	return nil
}

func (repo BookRepo) DeleteOne(bookId string) error {
	exist, err := repo.Exists(bookId)
	if err != nil {
		return logTrace(err)
	}

	if !exist {
		return nil
	}

	err = repo.db.Where("id = ?", bookId).Delete(&database.Book{}).Error
	if err != nil {
		return logTrace(err)
	}
	return nil
}

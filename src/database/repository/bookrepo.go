package repository

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/logging"
	"gorm.io/gorm"
)

type BookRepo struct {
	db *gorm.DB
}

func NewBookRepo(db *gorm.DB) *BookRepo {
	return &BookRepo{db: db}
}

func (repo BookRepo) FindOneByName(bookName string) (*database.Book, error) {
	res := database.Book{}
	err := repo.db.Where("name = ?", bookName).First(&res).Error
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	return &res, nil
}

func (repo BookRepo) Exists(bookId string) (bool, error) {
	var exists bool
	err := repo.db.Model(&database.Book{}).
		Select("count(*) > 0").
		Where("id = ?", bookId).
		Find(&exists).Error

	if err != nil {
		return false, logging.LogTrace(err)
	}
	return exists, nil
}

func (repo BookRepo) ExistsByName(bookName string) (bool, error) {
	var exists bool
	err := repo.db.Model(&database.Book{}).
		Select("count(*) > 0").
		Where("name = ?", bookName).
		Find(&exists).Error

	if err != nil {
		return false, logging.LogTrace(err)
	}
	return exists, nil
}

func (repo BookRepo) ExistAndFetchByName(bookName string) (*database.Book, error) {
	exist, err := repo.ExistsByName(bookName)
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	if !exist {
		return nil, logging.LogTrace(fmt.Errorf("cannot fetch book [%v] that does not exist", bookName))
	}

	book, err := repo.FindOneByName(bookName)
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	return book, nil
}

func (repo BookRepo) CreateOne(book *database.Book) error {
	if err := repo.db.Create(book).Error; err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

func (repo BookRepo) DeleteOne(bookId string) error {
	exist, err := repo.Exists(bookId)
	if err != nil {
		return logging.LogTrace(err)
	}

	if !exist {
		return nil
	}

	err = repo.db.Where("id = ?", bookId).Delete(&database.Book{}).Error
	if err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

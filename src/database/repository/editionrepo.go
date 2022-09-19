package repository

import (
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/logging"
	"gorm.io/gorm"
)

type EditionRepo struct {
	db *gorm.DB
}

func NewEditionRepo(db *gorm.DB) *EditionRepo {
	return &EditionRepo{db: db}
}

func (repo EditionRepo) Exists(EditionId string) (bool, error) {
	var exists bool
	err := repo.db.Model(&database.Edition{}).
		Select("count(*) > 0").
		Where("id = ?", EditionId).
		Find(&exists).Error

	if err != nil {
		return false, logging.LogTrace(err)
	}
	return exists, nil
}

func (repo EditionRepo) ExistByBookIdAndEditionNumber(bookId string, editionNumber int) (bool, error) {
	var exist bool
	err := repo.db.Model(&database.Edition{}).
		Select("count(*) > 0").
		Where("book_id = ? and edition_number = ?", bookId, editionNumber).
		Find(&exist).Error
	if err != nil {
		return false, logging.LogTrace(err)
	}
	return exist, nil
}

func (repo EditionRepo) FindByBookIdAndEditionNumber(bookId string, editionNumber int) (*database.Edition, error) {
	res := &database.Edition{}
	err := repo.db.Where("book_id = ? and edition_number = ?", bookId, editionNumber).Find(res).Error

	if err != nil {
		return nil, logging.LogTrace(err)
	}
	return res, nil
}

func (repo EditionRepo) FindNextEditionNumber(bookId string) (int, error) {
	var exists bool
	err := repo.db.Model(&database.Edition{}).
		Select("count(*) > 0").
		Where("book_id = ?", bookId).
		Find(&exists).Error

	if err != nil {
		return 0, logging.LogTrace(err)
	}
	if !exists {
		return 0, nil
	} else {
		var num = 0
		err := repo.db.Model(&database.Edition{}).Select("edition_number").Where("book_id = ?", bookId).Order("edition_number desc").Limit(1).Find(&num).Error
		if err != nil {
			return 0, logging.LogTrace(err)
		}
		return num + 1, nil
	}

}

func (repo EditionRepo) CreateOne(Edition *database.Edition) error {
	if err := repo.db.Create(Edition).Error; err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

func (repo EditionRepo) DeleteOne(EditionId string) error {
	exist, err := repo.Exists(EditionId)
	if err != nil {
		return logging.LogTrace(err)
	}

	if !exist {
		return nil
	}

	err = repo.db.Where("id = ?", EditionId).Delete(&database.Edition{}).Error
	if err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

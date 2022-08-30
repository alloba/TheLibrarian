package repository

import (
	"github.com/alloba/TheLibrarian/database"
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
		return false, logTrace(err)
	}
	return exists, nil
}

func (repo EditionRepo) FindNextEditionNumber(bookId string) (int, error) {
	var exists bool
	err := repo.db.Model(&database.Edition{}).
		Select("count(*) > 0").
		Where("book_id = ?", bookId).
		Find(&exists).Error

	if err != nil {
		return 0, logTrace(err)
	}
	if !exists {
		return 0, nil
	} else {
		var num = 0
		err := repo.db.Model(&database.Edition{}).Select("edition_number").Where("book_id = ?", bookId).Order("edition_number desc").Limit(1).Find(&num).Error
		if err != nil {
			return 0, logTrace(err)
		}
		return num + 1, nil
	}

}

func (repo EditionRepo) CreateOne(Edition *database.Edition) error {
	if err := repo.db.Create(Edition).Error; err != nil {
		return logTrace(err)
	}
	return nil
}

func (repo EditionRepo) DeleteOne(EditionId string) error {
	exist, err := repo.Exists(EditionId)
	if err != nil {
		return logTrace(err)
	}

	if !exist {
		return nil
	}

	err = repo.db.Where("id = ?", EditionId).Delete(&database.Edition{}).Error
	if err != nil {
		return logTrace(err)
	}
	return nil
}

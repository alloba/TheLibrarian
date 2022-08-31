package repository

import (
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
)

type ChapterRepo struct {
	db *gorm.DB
}

func NewChapterRepo(db *gorm.DB) *ChapterRepo {
	return &ChapterRepo{db: db}
}

func (repo ChapterRepo) FindOneByNameAndEdition(ChapterName string, editionId int) (*database.Chapter, error) {
	res := database.Chapter{}
	err := repo.db.Where("name = ? and edition_id = ?", ChapterName, editionId).First(&res).Error
	if err != nil {
		return nil, logTrace(err)
	}
	return &res, nil
}

func (repo ChapterRepo) FindAllByEditionId(editionId string) (*[]database.Chapter, error) {
	res := make([]database.Chapter, 0)
	err := repo.db.Where("edition_id = ?", editionId).Find(&res).Error
	if err != nil {
		return nil, logTrace(err)
	}
	return &res, nil
}

func (repo ChapterRepo) Exists(ChapterId string) (bool, error) {
	var exists bool
	err := repo.db.Model(&database.Chapter{}).
		Select("count(*) > 0").
		Where("id = ?", ChapterId).
		Find(&exists).Error

	if err != nil {
		return false, logTrace(err)
	}
	return exists, nil
}

func (repo ChapterRepo) ExistsByNameAndEdition(ChapterName string, editionId int) (bool, error) {
	var exists bool
	err := repo.db.Model(&database.Chapter{}).
		Select("count(*) > 0").
		Where("name = ? and edition_id = ?", ChapterName, editionId).
		Find(&exists).Error

	if err != nil {
		return false, logTrace(err)
	}
	return exists, nil
}

func (repo ChapterRepo) CreateOne(Chapter *database.Chapter) error {
	if err := repo.db.Create(Chapter).Error; err != nil {
		return logTrace(err)
	}
	return nil
}

func (repo ChapterRepo) DeleteOne(ChapterId string) error {
	exist, err := repo.Exists(ChapterId)
	if err != nil {
		return logTrace(err)
	}

	if !exist {
		return nil
	}

	err = repo.db.Where("id = ?", ChapterId).Delete(&database.Chapter{}).Error
	if err != nil {
		return logTrace(err)
	}
	return nil
}

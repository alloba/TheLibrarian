package repository

import (
	"gorm.io/gorm"
)

type RepoManager struct {
	Record  *RecordRepo
	Book    *BookRepo
	Edition *EditionRepo
	Page    *PageRepo
	Chapter *ChapterRepo
	Db      *gorm.DB
}

func NewRepoManager(db *gorm.DB) *RepoManager {
	return &RepoManager{
		Record:  NewRecordRepo(db),
		Book:    NewBookRepo(db),
		Edition: NewEditionRepo(db),
		Page:    NewPageRepo(db),
		Chapter: NewChapterRepo(db),
		Db:      db,
	}
}

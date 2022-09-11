package repository

import (
	"gorm.io/gorm"
)
//TODO: This honestly feels like a progressively worse idea the longer i use it.
//      Either i should fall into full dependency injection, 
//      or i should explicitly insert pointers in services that need them...
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

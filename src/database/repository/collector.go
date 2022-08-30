package repository

import (
	"gorm.io/gorm"
)

type RepoManager struct {
	Book    *BookRepo
	Page    *PageRepo
	Edition *EditionRepo
	Record  *RecordRepo
}

func NewRepoManager(db *gorm.DB) *RepoManager {
	return &RepoManager{
		Book:    NewBookRepo(db),
		Page:    NewPageRepo(db),
		Edition: NewEditionRepo(db),
		Record:  NewRecordRepo(db),
	}
}

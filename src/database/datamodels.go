package database

import (
	"time"
)

//Tabler provides a mechanism to override gorm table name defaults
type Tabler interface {
	TableName() string
}

func (Record) TableName() string {
	return "record"
}

type Record struct {
	Hash             string    `gorm:"primaryKey"`
	FilePointer      string    `gorm:"not null"`
	Name             string    `gorm:"not null"`
	Extension        string    `gorm:"not null"`
	DateFileModified time.Time `gorm:"not null"`
	DateCreated      time.Time `gorm:"not null"`
	DateModified     time.Time `gorm:"not null"`
}

func (Book) TableName() string {
	return "book"
}

type Book struct {
	Uuid         string    `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
	DateCreated  time.Time `gorm:"not null"`
	DateModified time.Time `gorm:"not null"`
}

type Page struct {
	Uuid         string
	RecordHash   string
	BookUuid     string
	EditionUuid  string
	DateCreated  time.Time
	DateModified time.Time
}

type Edition struct {
	Uuid          string
	EditionNumber string
	BookUuid      string
	DateCreated   time.Time
	DateModified  time.Time
}

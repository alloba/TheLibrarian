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

func (Edition) TableName() string {
	return "edition"
}

type Edition struct {
	Uuid          string    `gorm:"primaryKey"`
	EditionNumber int       `gorm:"not null"`
	BookUuid      string    `gorm:"not null" gorm:"references book.uuid"`
	Book          *Book     `gorm:"not null" gorm:"foreignKey:BookUuid"`
	DateCreated   time.Time `gorm:"not null"`
	DateModified  time.Time `gorm:"not null"`
}

type Page struct {
	Uuid         string
	RecordHash   string
	BookUuid     string
	EditionUuid  string
	DateCreated  time.Time
	DateModified time.Time
}

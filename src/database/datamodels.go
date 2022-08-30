package database

import (
	"time"
)

//Tabler provides a mechanism to override gorm table name defaults
//type Tabler interface {
//	TableName() string
//}

//func (Record) TableName() string {
//	return "record"
//}

type Record struct {
	Id               string    `gorm:"primaryKey" gorm:"constraint:OnDelete:CASCADE;"`
	FilePointer      string    `gorm:"not null"`
	Name             string    `gorm:"not null"`
	Extension        string    `gorm:"not null"`
	DateFileModified time.Time `gorm:"not null"`
	DateCreated      time.Time `gorm:"not null"`
	DateModified     time.Time `gorm:"not null"`
}

//func (Book) TableName() string {
//	return "book"
//}

type Book struct {
	Id           string    `gorm:"primaryKey" gorm:"constraint:OnDelete:CASCADE;"`
	Name         string    `gorm:"not null"`
	DateCreated  time.Time `gorm:"not null"`
	DateModified time.Time `gorm:"not null"`
}

//func (Edition) TableName() string {
//	return "edition"
//}

type Edition struct {
	Id            string    `gorm:"primaryKey"`
	EditionNumber int       `gorm:"not null"`
	BookId        string    `gorm:"not null" gorm:"constraint:OnDelete:CASCADE;"`
	Book          *Book     `gorm:"not null" gorm:"constraint:OnDelete:CASCADE;"`
	Pages         *[]Page   `gorm:"foreignKey:EditionId" gorm:"constraint:OnDelete:CASCADE;"`
	DateCreated   time.Time `gorm:"not null"`
	DateModified  time.Time `gorm:"not null"`
}

//func (Page) TableName() string {
//	return "page"
//}

type Page struct {
	Id           string    `gorm:"primaryKey"`
	RecordId     string    `gorm:"not null"`
	Record       Record    `gorm:"not null" gorm:"foreignKey:RecordId" gorm:"constraint:OnDelete:CASCADE;"`
	EditionId    string    `gorm:"not null" gorm:"constraint:OnDelete:CASCADE;"`
	Edition      *Edition  `gorm:"not null" gorm:"foreignKey:EditionId" gorm:"constraint:OnDelete:CASCADE;"`
	DateCreated  time.Time `gorm:"not null"`
	DateModified time.Time `gorm:"not null"`
}

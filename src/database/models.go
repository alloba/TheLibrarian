package database

import (
	"time"
)

type Record struct {
	ID               string    `gorm:"column:id"                 gorm:"primaryKey"` // gorm:"constraint:OnDelete:CASCADE;"`
	FilePointer      string    `gorm:"column:file_pointer"       gorm:"not null"`
	Name             string    `gorm:"column:name"               gorm:"not null"`
	Extension        string    `gorm:"column:extension"          gorm:"not null"`
	DateFileModified time.Time `gorm:"column:date_file_modified" gorm:"not null"`
	DateCreated      time.Time `gorm:"column:date_created"       gorm:"not null"`
	DateModified     time.Time `gorm:"column:date_modified"      gorm:"not null"`
}

type Book struct {
	ID           string    `gorm:"column:id"            gorm:"primaryKey"`
	Name         string    `gorm:"column:name"          gorm:"not null"`
	DateCreated  time.Time `gorm:"column:date_created"  gorm:"not null"`
	DateModified time.Time `gorm:"column:date_modified" gorm:"not null"`
}

type Edition struct {
	ID            string    `gorm:"column:id"             gorm:"primaryKey"`
	EditionNumber int       `gorm:"column:edition_number" gorm:"not null"`
	BookID        string    `gorm:"column:book_id"        gorm:"not null"`
	Book          *Book     `gorm:"foreignKey:ID"`
	Pages         *[]Page   `gorm:"foreignKey:ID"`
	DateCreated   time.Time `gorm:"column:date_created"   gorm:"not null"`
	DateModified  time.Time `gorm:"column:date_modified"  gorm:"not null"`
}

type Page struct {
	ID           string    `gorm:"column:id"            gorm:"primaryKey"`
	RecordID     string    `gorm:"column:record_id"     gorm:"not null"`
	Record       Record    `gorm:"foreignKey:ID"`
	EditionID    string    `gorm:"column:edition_id"    gorm:"not null"`
	Edition      *Edition  `gorm:"foreignKey:ID"`
	DateCreated  time.Time `gorm:"column:date_created"  gorm:"not null"`
	DateModified time.Time `gorm:"column:date_modified" gorm:"not null"`
}

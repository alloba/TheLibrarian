package database

import (
	"time"
)

type Record struct {
	ID               string    `gorm:"column:id"                 gorm:"primaryKey"`
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
	Name          string    `gorm:"column:name"           gorm:"not null"`
    EditionNumber int       `gorm:"column:edition_number" gorm:"not null"`
	BookID        string    `gorm:"column:book_id"        gorm:"not null"`
	DateCreated   time.Time `gorm:"column:date_created"   gorm:"not null"`
	DateModified  time.Time `gorm:"column:date_modified"  gorm:"not null"`
}

type Page struct {
	ID           string    `gorm:"column:id"            gorm:"primaryKey"`
	RecordID     string    `gorm:"column:record_id"     gorm:"not null"`
    EditionID    string    `gorm:"column:edition_id"    gorm:"not null"`
	RelativePath string    `gorm:"column:relative_path" gorm:"not null"`
	DateCreated  time.Time `gorm:"column:date_created"  gorm:"not null"`
	DateModified time.Time `gorm:"column:date_modified" gorm:"not null"`
}

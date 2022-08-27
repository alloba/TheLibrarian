package database

import "time"

type Record struct {
	Hash             string
	FilePointer      string
	Name             string
	Extension        string
	DateFileModified time.Time
	DateCreated      time.Time
	DateModified     time.Time
}

type Page struct {
	Uuid         string
	RecordHash   string
	BookUuid     string
	EditionUuid  string
	DateCreated  time.Time
	DateModified time.Time
}

type Book struct {
	Uuid         string
	Name         string
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

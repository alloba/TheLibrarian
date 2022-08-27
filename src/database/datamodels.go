package database

import "time"

type Record struct {
	hash             string
	filePointer      string
	name             string
	extension        string
	dateFileModified time.Time

	dateCreated  time.Time
	dateModified time.Time
}

type Page struct {
	uuid        string
	recordHash  string
	bookUuid    string
	editionUuid string

	dateCreated  time.Time
	dateModified time.Time
}

type Book struct {
	uuid string
	name string

	dateCreated  time.Time
	dateModified time.Time
}

type Edition struct {
	uuid          string
	editionNumber string
	bookUuid      string

	dateCreated  time.Time
	dateModified time.Time
}

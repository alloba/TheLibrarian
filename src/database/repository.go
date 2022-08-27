package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func NewRecordRepo(db *sql.DB) RecordRepo {
	return RecordRepo{
		db:      db,
		FindAll: getAllRecords(db),
	}
}

type RecordRepo struct {
	db      *sql.DB
	FindAll findAllRecords
}
type findAllRecords func() ([]Record, error)

func getAllRecords(db *sql.DB) func() ([]Record, error) {
	return func() ([]Record, error) {
		log.Printf("beginning operation getAllRecords.")

		stmt, err := db.Prepare(`select 
        hash, 
        file_pointer, 
        name, 
        extension, 
        date_created, 
        date_modified, 
        date_file_modified
        from record
       `)
		if err != nil {
			return nil, fmt.Errorf("could not form statement for getAllRecords: %w", err)
		}
		defer stmt.Close()

		rows, err := stmt.Query()
		if err != nil {
			return nil, fmt.Errorf("could not execute statement for getAllRecords: %w", err)
		}
		defer rows.Close()

		var hash string
		var filePointer string
		var name string
		var extension string
		var dateFileModified time.Time
		var dateCreated time.Time
		var dateModified time.Time
		var records = make([]Record, 0)

		for rows.Next() {
			err := rows.Scan(&hash, &filePointer, &name, &extension, &dateFileModified, &dateCreated, &dateModified)
			if err != nil {
				return nil, fmt.Errorf("could not poll row results for getAllRecords: %w", err)
			}
			records = append(records, Record{
				hash:             hash,
				filePointer:      filePointer,
				name:             name,
				extension:        extension,
				dateFileModified: dateFileModified,
				dateCreated:      dateCreated,
				dateModified:     dateModified,
			})
		}
		log.Printf("completed operation getAllRecords. returning %v results", len(records))
		return records, nil
	}
}

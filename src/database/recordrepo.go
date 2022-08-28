package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func NewRecordRepo(db *sql.DB) RecordRepo {
	return RecordRepo{
		db:         db,
		FindAll:    getAllRecords(db),
		SaveOne:    saveSingleRecord(db),
		FindByHash: findSingleRecordByHash(db),
	}
}

type RecordRepo struct {
	db         *sql.DB
	FindAll    findAllRecords
	SaveOne    saveRecord
	FindByHash findRecordByHash
}
type findAllRecords func() ([]Record, error)
type saveRecord func(record *Record) error
type findRecordByHash func(hash string) (*Record, error)

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
        from record`)

		if err != nil {
			return nil, fmt.Errorf("could not form statement for getAllRecords: %w", err)
		}
		defer stmt.Close()

		rows, err := stmt.Query()
		if err != nil {
			return nil, fmt.Errorf("could not execute statement for getAllRecords: %w", err)
		}
		defer rows.Close()

		var records = make([]Record, 0)
		for rows.Next() {
			record, err := createRecordFromDbResult(rows)
			if err != nil {
				return nil, fmt.Errorf("could not poll row results for getAllRecords: %w", err)
			}
			records = append(records, *record)
		}
		log.Printf("completed operation getAllRecords. returning %v results", len(records))
		return records, nil
	}
}

func findSingleRecordByHash(db *sql.DB) func(hash string) (*Record, error) {
	return func(hash string) (*Record, error) {
		log.Printf("finding record by hash %v", hash)
		stmt, err := db.Prepare(`select r.hash, r.file_pointer, r.name, r.extension, r.date_created, r.date_modified, r.date_file_modified from record r where r.hash = ?`)
		if err != nil {
			return nil, fmt.Errorf("could not form statement for findRecordByHash: %w", err)
		}
		defer stmt.Close()

		res, err := stmt.Query(hash)
		defer res.Close()
		if err != nil {
			return nil, fmt.Errorf("could not execute findRecordByHash: %w", err)
		}

		record, err := createRecordFromDbResult(res)
		if err != nil {
			return nil, fmt.Errorf("could not findRecordByHash: %w", err)
		}

		log.Printf("completed finding record by hash %v.", record)
		return record, nil
	}
}

func createRecordFromDbResult(result *sql.Rows) (*Record, error) {
	var hash string
	var filePointer string
	var name string
	var extension string
	var dateFileModified time.Time
	var dateCreated time.Time
	var dateModified time.Time

	var nRow = result.Next()
	if nRow {
		var err = result.Scan(
			&hash,
			&filePointer,
			&name,
			&extension,
			&dateFileModified,
			&dateCreated,
			&dateModified,
		)
		if err != nil {
			return nil, fmt.Errorf("error reading row from database: %w", err)
		}
		return &Record{
			Hash:             hash,
			FilePointer:      filePointer,
			Name:             name,
			Extension:        extension,
			DateFileModified: dateFileModified,
			DateCreated:      dateCreated,
			DateModified:     dateModified,
		}, nil
	} else {
		return nil, fmt.Errorf("unable to create Record from db result - no remaining rows")
	}

}

func saveSingleRecord(db *sql.DB) func(record *Record) error {
	return func(record *Record) error {
		log.Printf("saving record to database")
		stmt, err := db.Prepare(`insert into record(hash, file_pointer, name, extension, date_created, date_file_modified, date_modified) values (?,?,?,?,?,?,?)`)
		if err != nil {
			return fmt.Errorf("unable to form statement for saveSingleRecord: %w", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(record.Hash, record.FilePointer, record.Name, record.Extension, record.DateCreated, record.DateFileModified, record.DateModified)
		if err != nil {
			log.Printf("Failed to save record with hash %v to the database. Reason: %v", record.Hash, err.Error())
			return fmt.Errorf("unable to save record to database: %w", err)
		}
		log.Printf("completed saving record to database")
		return nil
	}
}

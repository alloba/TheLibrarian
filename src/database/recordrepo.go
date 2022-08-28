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

type recordSqlObject struct {
	hash             sql.NullString
	filePointer      sql.NullString
	name             sql.NullString
	extension        sql.NullString
	dateFileModified sql.NullTime
	dateCreated      sql.NullTime
	dateModified     sql.NullTime
}

type findAllRecords func() ([]Record, error)
type saveRecord func(record *Record) error
type findRecordByHash func(hash string) (*Record, error)

func getAllRecords(db *sql.DB) func() ([]Record, error) {
	statement, err := db.Prepare(`select hash, file_pointer, name, extension, date_created, date_modified, date_file_modified from record`)
	if err != nil {
		log.Fatalf("could not create prepared statement for getAllRecords - %v", err.Error())
	}

	return func() ([]Record, error) {
		log.Printf("beginning operation getAllRecords.")

		if err != nil {
			return nil, fmt.Errorf("could not form statement for getAllRecords: %w", err)
		}

		rows, err := statement.Query()
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
	statement, err := db.Prepare(`select r.hash, r.file_pointer, r.name, r.extension, r.date_created, r.date_modified, r.date_file_modified from record r where r.hash = ?`)
	if err != nil {
		log.Fatalf("could not create prepared statement for findSingleRecordByHash - %v", err.Error())
	}

	return func(hash string) (*Record, error) {
		log.Printf("finding record by hash %v", hash)

		res, err := statement.Query(hash)
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

func saveSingleRecord(db *sql.DB) func(record *Record) error {
	statement, err := db.Prepare(`insert into record(hash, file_pointer, name, extension, date_created, date_file_modified, date_modified) values (?,?,?,?,?,?,?)`)
	if err != nil {
		log.Fatalf("cold not create prepared statement for saveSingleRecord - %v", err.Error())
	}

	return func(record *Record) error {
		log.Printf("saving record to database")
		_, err = statement.Exec(record.Hash, record.FilePointer, record.Name, record.Extension, record.DateCreated, record.DateFileModified, record.DateModified)
		if err != nil {
			log.Printf("Failed to save record with hash %v to the database. Reason: %v", record.Hash, err.Error())
			return fmt.Errorf("unable to save record to database: %w", err)
		}
		log.Printf("completed saving record to database")
		return nil
	}
}

// given a list of rows from the database, attempt to scan the next row and form a record object
// returns an error if there are no remaining rows to scan,
// or if there was some issue while forming a record object.
func createRecordFromDbResult(result *sql.Rows) (*Record, error) {
	var holder = recordSqlObject{}
	var nRow = result.Next()
	if nRow {
		var err = result.Scan(&holder.hash, &holder.filePointer, &holder.name, &holder.extension, &holder.dateFileModified, &holder.dateCreated, &holder.dateModified)
		if err != nil {
			return nil, fmt.Errorf("error reading row from database: %w", err)
		}
		finalRecord, err := convertRecordSqlToRecord(&holder)
		if err != nil {
			return nil, fmt.Errorf("unable to convert sql object to record - %v", err.Error())
		}

		return finalRecord, nil
	} else {
		return nil, fmt.Errorf("unable to create Record from db result - no remaining rows")
	}
}

// helper function for converting between a recordSqlObject and a final Record object.
// this allows a place to handle null checks for data coming from the db table.
// for the current schema, this means throwing errors on null values. but in the future there may be some
// allowances for null fields. This provides a template going forward.
func convertRecordSqlToRecord(object *recordSqlObject) (*Record, error) {
	var out = Record{
		Hash:             "",
		FilePointer:      "",
		Name:             "",
		Extension:        "",
		DateFileModified: time.Time{},
		DateCreated:      time.Time{},
		DateModified:     time.Time{},
	}

	if object.hash.Valid {
		out.Hash = object.hash.String
	} else {
		return nil, fmt.Errorf("null value encountered when creating record object for field 'hash'")
	}

	if object.filePointer.Valid {
		out.FilePointer = object.filePointer.String
	} else {
		return nil, fmt.Errorf("null value encountered when creating record object for field 'filePointer'")
	}

	if object.name.Valid {
		out.Name = object.name.String
	} else {
		return nil, fmt.Errorf("null value encountered when creating record object for field 'name'")
	}

	if object.extension.Valid {
		out.Extension = object.extension.String
	} else {
		return nil, fmt.Errorf("null value encountered when creating record object for field 'extension'")
	}

	if object.dateFileModified.Valid {
		out.DateFileModified = object.dateFileModified.Time
	} else {
		return nil, fmt.Errorf("null value encountered when creating record object for field 'dateFileModified'")
	}

	if object.dateCreated.Valid {
		out.DateCreated = object.dateCreated.Time
	} else {
		return nil, fmt.Errorf("null value encountered when creating record object for field 'dateCreated'")
	}

	if object.dateModified.Valid {
		out.DateModified = object.dateModified.Time
	} else {
		return nil, fmt.Errorf("null value encountered when creating record object for field 'dateModified'")
	}

	return &out, nil
}

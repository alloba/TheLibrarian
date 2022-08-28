package database

import (
	"database/sql"
	"log"
	"testing"
	"time"
)

func Test_SaveAndEchoSingleRecord(t *testing.T) {
	var db = Connect("../../schema/library.db")
	defer db.Close()

	var recordRepo = NewRecordRepo(db)
	deleteTestRecords(db)
	defer deleteTestRecords(db)

	var record = Record{
		Hash:             "testhash1",
		FilePointer:      "filepoint",
		Name:             "filename",
		Extension:        "fileextensionj",
		DateFileModified: time.Now(),
		DateCreated:      time.Now(),
		DateModified:     time.Now(),
	}
	t.Run("singleRecord", func(t *testing.T) {
		var err = recordRepo.SaveOne(&record)
		if err != nil {
			t.Errorf("Failed test %v while saving record - %v", t.Name(), err.Error())
		}

		record, err := recordRepo.FindByHash("testhash1")
		if err != nil {
			t.Errorf("failed test %v while searching for record - %v", t.Name(), err.Error())
		}

		if record.Hash != "testhash1" {
			t.Errorf("failed test %v. returned record hash does not match - %v", t.Name(), err.Error())
		}
	})
}

func deleteTestRecords(db *sql.DB) {
	_, err := db.Exec("delete from record where lower(hash) like 'test%'")
	if err != nil {
		log.Fatalf("failed to delete all records from database for testing: %v", err.Error())
	}
}

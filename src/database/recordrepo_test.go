package database

import (
	"gorm.io/gorm"
	"testing"
	"time"
)

func Test_SaveAndEchoSingleRecord(t *testing.T) {
	//todo need to have a separate testing database away from the actual in-use one.
	//todo i really need a central way to manage database location.
	var db = Connect("../../out/library.db")

	var recordRepo = NewRecordRepo(db)
	deleteTestRecords(db)
	defer deleteTestRecords(db)

	var record = Record{
		Hash:             "testhash1",
		FilePointer:      "filepoint",
		Name:             "filenameeeeeeeeeez",
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

		r, err := recordRepo.FindByHash("testhash1")
		if err != nil {
			t.Errorf("failed test %v while searching for record - %v", t.Name(), err.Error())
		}

		if r.Hash != "testhash1" {
			t.Errorf("failed test %v. returned record hash does not match - %v", t.Name(), err.Error())
		}
	})
}

func Test_RecordDoesNotExist(t *testing.T) {
	var db = Connect("../../out/library.db")

	var recordRepo = NewRecordRepo(db)
	deleteTestRecords(db)
	defer deleteTestRecords(db)

	t.Run("noRecord", func(t *testing.T) {
		_, err := recordRepo.FindByHash("test_zzzzz")
		if err == nil {
			t.Errorf("failed test %v while searching for record - %v", t.Name(), err.Error())
		}
	})
}

func Test_getAllRecords(t *testing.T) {
	var db = Connect("../../out/library.db")
	var recordRepo = NewRecordRepo(db)
	deleteTestRecords(db)
	defer deleteTestRecords(db)
	var record = getDummyRecord()
	recordRepo.SaveOne(record)

	t.Run("getAllRecords", func(t *testing.T) {
		records, err := recordRepo.FindAll()
		if err != nil || len(records) == 0 {
			t.Errorf("failed to findAll records - %v", err.Error())
		}
	})

}

func deleteTestRecords(db *gorm.DB) {
	db.Where("hash like ?", "test%").Delete(&Record{})
}

func getDummyRecord() *Record {
	return &Record{
		Hash:             "testhash1",
		FilePointer:      "filepoint",
		Name:             "filename",
		Extension:        "fileextensionj",
		DateFileModified: time.Now(),
		DateCreated:      time.Now(),
		DateModified:     time.Now(),
	}
}

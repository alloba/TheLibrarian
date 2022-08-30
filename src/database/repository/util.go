package repository

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"runtime"
	"strings"
	"time"
)

var integrationTestDbPath = "../../../out/library_integration_test.db"

func logTrace(err error) error {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("ERR_CANNOT_TRACE_CALLER: %v", err.Error())
	}

	fullName := fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	methodName := strings.Split(fullName, "/")[len(strings.Split(fullName, "/"))-1]

	return fmt.Errorf("%v: %v", methodName, err.Error())
}

func getTestRecord(id string) *database.Record {
	return &database.Record{
		ID:               id,
		FilePointer:      "testpointer",
		Name:             "testname",
		Extension:        "test",
		DateFileModified: time.Now(),
		DateCreated:      time.Now(),
		DateModified:     time.Now(),
	}
}

func getTestPage(id string, editionId string) *database.Page {
	record := getTestRecord("testRecordAssociatedWithPage" + id)
	//edition := getTestEdition("testEditionAssociatedWithPage" + id)

	return &database.Page{
		ID:        id,
		RecordID:  record.ID,
		Record:    *record,
		EditionID: editionId,
		//Edition:      edition,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}
}

func getTestBook(id string) *database.Book {
	return &database.Book{
		ID:           id,
		Name:         "testname",
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}
}

package repository

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
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

func getTestBook(id string) *database.Book {
	return &database.Book{
		ID:           id,
		Name:         "testname",
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}
}

func WipeTestDatabase(db *gorm.DB) {
	//delete page
	//delete edition
	//delete book
	//delete record
	db.Where("1=1").Delete(&database.Page{})
	db.Where("1=1").Delete(&database.Chapter{})
	db.Where("1=1").Delete(&database.Edition{})
	db.Where("1=1").Delete(&database.Book{})
	db.Where("1=1").Delete(&database.Record{})
}

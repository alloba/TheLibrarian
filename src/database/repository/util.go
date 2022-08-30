package repository

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
	"runtime"
	"strings"
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

func cleanUpTestRecords(db *gorm.DB) {
	err := db.Where("id like ?", "test%").Delete(&database.Record{}).Error
	if err != nil {
		panic(logTrace(err))
	}
}

func cleanUpTestBooks(db *gorm.DB) {
	err := db.Where("id like ?", "test%").Delete(&database.Book{}).Error
	if err != nil {
		panic(logTrace(err))
	}
}

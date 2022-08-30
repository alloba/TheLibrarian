package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"testing"
)

func Test_createRecordFromFile(t *testing.T) {
	var db = database.Connect("../../out/library.db")
	var repo = database.NewRecordRepo(db)
	var basepath = "./"
	var fileservice = NewFileService(basepath)
	var service = NewRecordService(repo, fileservice)
	t.Run("createRecord", func(t *testing.T) {
		rec, err := service.CreateRecordData("./recordservice_test.go")
		if err != nil {
			t.Errorf("failed to create record - %v", err.Error())
		}
		fmt.Printf("%#v", rec)
	})
}

func Test_persistNewRecord(t *testing.T) {
	var db = database.Connect("../../out/library.db")
	var repo = database.NewRecordRepo(db)
	var basepath = "../../out/filebin"
	var fileservice = NewFileService(basepath)
	var service = NewRecordService(repo, fileservice)
	t.Run("createRecord", func(t *testing.T) {
		rec, err := service.CreateRecordData("./recordservice_test.go")
		if err != nil {
			t.Errorf("failed to create record - %v", err.Error())
		}
		effect, err := service.PersistRecordIfUnique(&RecordZip{
			OriginPath: "./recordservice_test.go",
			RecordItem: rec,
		})
		if err != nil {
			t.Errorf("couldnt persist record - %v", err.Error())
		}
		if !effect {
			fmt.Printf("no action taken, record already exists - %v", rec.Id)
		}
		if effect {
			fmt.Printf("file pushed to database and library bin - %v", rec.FilePointer)
		}

		fmt.Printf("%#v", rec)
	})
}

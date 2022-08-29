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
	var service = New(repo, basepath)
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
	var service = New(repo, basepath)
	t.Run("createRecord", func(t *testing.T) {
		rec, err := service.CreateRecordData("./recordservice_test.go")
		if err != nil {
			t.Errorf("failed to create record - %v", err.Error())
		}
		effect, err := service.PersistRecordIfUnique(rec, "./recordservice_test.go") //todo not a huge fan of having to supply the file location a second time here.
		if err != nil {
			t.Errorf("couldnt persist record - %v", err.Error())
		}
		if !effect {
			fmt.Printf("no action taken, record already exists - %v", rec.Hash)
		}
		if effect {
			fmt.Printf("file pushed to database and library bin - %v", rec.FilePointer)
		}

		fmt.Printf("%#v", rec)
	})
}

func Test_generateUniqueSubPath(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"case1", args{hash: "asdfhash"}},
		{"case2", args{hash: "^&*hkjalsdfhkjl__!!jfjfjfjfj.zip"}},
		{"case2", args{hash: "a.a.a."}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generateUniqueSubPath(tt.args.hash)
		})
	}
}

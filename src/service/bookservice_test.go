package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"gorm.io/gorm"
	"testing"
)

func TestBookService_createNewBook(t *testing.T) {
	var db = database.Connect("../../out/library.db")
	var recordRepo = database.NewRecordRepo(db)
	var bookRepo = database.NewBookRepo(db)
	var editionRepo = database.NewEditionRepo(db)

	var fileService = NewFileService("../../out/filebin")
	var recordService = NewRecordService(recordRepo, fileService)

	deleteTestEditions(db)
	//defer deleteTestEditions(db)

	var bookService = NewBookService(recordService, fileService, bookRepo, editionRepo)
	fmt.Printf("%v", bookService)
	t.Run("createBookFromDirectory", func(t *testing.T) {
		res1, err1 := bookService.CreateNewBook("testBook1", "../")
		if err1 != nil {
			t.Fatalf("could not create and save book edition - %v", err1.Error())
		}
		res2, err2 := bookService.CreateNewBook("testBook2", "../")
		if err2 != nil {
			t.Fatalf("could not create and save book edition - %v", err2.Error())
		}
		res3, err3 := bookService.CreateNewBook("testBook3", "../")
		if err3 != nil {
			t.Fatalf("could not create and save book edition - %v", err3.Error())
		}
		t.Logf("%#v", res1)
		t.Logf("%#v", res2)
		t.Logf("%#v", res3)
	})
}

func deleteTestEditions(db *gorm.DB) {
	db.Where("1=1").Delete(&database.Page{})
	db.Where("1=1").Delete(&database.Edition{})
	db.Where("1=1").Delete(&database.Book{})
	db.Where("1=1").Delete(&database.Record{})
	//db.Where("hash like ?", "test%").Delete(&Record{})
}

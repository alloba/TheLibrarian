package database

import (
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestBookRepo_FindOne(t *testing.T) {
	var db = Connect("../../out/library.db")
	var bookRepo = NewBookRepo(db)
	deleteTestBookEntries(db)
	defer deleteTestBookEntries(db)

	var testBook = Book{
		Id:           "testUuid1",
		Name:         "testBook1",
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}
	err := bookRepo.SaveOne(&testBook)
	if err != nil {
		t.Errorf("failed to initialize book FindOne test - could not save to database - %v", err.Error())
	}

	t.Run("findOneExists", func(t *testing.T) {
		res, err := bookRepo.FindOne("")
		if err != nil {
			t.Errorf("could not find book for uuid %v - %v", "", err.Error())
		}
		t.Log(res)
	})

	t.Run("findOneDoesNotExist", func(t *testing.T) {
		res, err := bookRepo.FindOne("aaa")
		if err == nil {
			t.Errorf("book returned for supposedly non-existent record %v", "")
		}
		t.Log(res)
	})
}

func TestBookRepo_SaveOne(t *testing.T) {
	var db = Connect("../../out/library.db")
	var bookRepo = NewBookRepo(db)
	deleteTestBookEntries(db)
	defer deleteTestBookEntries(db)

	t.Run("saveBookToDbNotExist", func(t *testing.T) {
		var testBook = Book{
			Id:           "testUuid1",
			Name:         "testBook1",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		}
		err := bookRepo.SaveOne(&testBook)
		if err != nil {
			t.Errorf("failed to save book to db %v - %v", testBook.Id, err.Error())
		}
	})

	t.Run("saveBookToDbExists", func(t *testing.T) {
		var testBook = Book{
			Id:           "testUuid1",
			Name:         "testBook1",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		}
		err := bookRepo.SaveOne(&testBook)
		if err == nil {
			t.Errorf("succeeded to save book to db when should be disallowed %v", testBook.Id)
		}
	})
}

func deleteTestBookEntries(db *gorm.DB) {
	db.Where("uuid like ?", "test%").Delete(&Book{})
}

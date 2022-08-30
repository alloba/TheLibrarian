package database

import (
	"gorm.io/gorm"
	"testing"
	"time"
)

//func TestEditionRepo_SaveOne(t *testing.T) {
//	var db = Connect("../../out/library.db")
//
//	var editionRepo = NewEditionRepo(db)
//	var bookRepo = NewBookRepo(db)
//
//	deleteTestEditionEntries(db)
//	deleteTestBookEntries(db)
//	defer deleteTestEditionEntries(db)
//	defer deleteTestBookEntries(db)
//
//	t.Run("saveEditionNew", func(t *testing.T) {
//		var book = Book{
//			Id:         "testBookuuid",
//			Name:         "testBook",
//			DateCreated:  time.Now(),
//			DateModified: time.Now(),
//		}
//		err := bookRepo.SaveOne(&book)
//		if err != nil {
//			t.Errorf("faled to prep saveEditionNew test, could not save book - %v", err.Error())
//		}
//
//		var edition = Edition{
//			Id:          "testEditionUuid",
//			EditionNumber: 0,
//			BookId:      "testBookuuid",
//			Book:          nil,
//			DateCreated:   time.Now(),
//			DateModified:  time.Now(),
//		}
//		err = editionRepo.SaveOne(&edition)
//		if err != nil {
//			t.Errorf("could not save edition for saveEditionNew test - %v", err.Error())
//		}
//	})
//
//	t.Run("saveEditionInvalidBookUuid", func(t *testing.T) {
//		var edition = Edition{
//			Id:          "testEditionUuidFailTarget",
//			EditionNumber: 0,
//			BookId:      "fakeuuid",
//			Book:          nil,
//			DateCreated:   time.Now(),
//			DateModified:  time.Now(),
//		}
//
//		err := editionRepo.SaveOne(&edition)
//		if err == nil {
//			t.Errorf("completed edition save when should not be possible due to fkey: %v", edition.Book.Id)
//		}
//	})
//}

func TestEditionRepo_SaveOne(t *testing.T) {
	var db = Connect("../../out/library.db")

	var editionRepo = NewEditionRepo(db)
	//var bookRepo = NewBookRepo(db)

	//deleteTestEditionEntries(db)
	//deleteTestBookEntries(db)
	////defer deleteTestBookEntries(db)
	////defer deleteTestEditionEntries(db)
	//
	//var book = Book{
	//	Id:         "testBookuuid",
	//	Name:         "testBook",
	//	DateCreated:  time.Now(),
	//	DateModified: time.Now(),
	//}
	//err := bookRepo.SaveOne(&book)
	//if err != nil {
	//	t.Errorf("faled to prep saveEditionNew test, could not save book - %v", err.Error())
	//}
	//
	//var edition = Edition{
	//	Id:          "testEditionUuid",
	//	EditionNumber: 0,
	//	BookId:      "testBookuuid",
	//	Book:          nil,
	//	DateCreated:   time.Now(),
	//	DateModified:  time.Now(),
	//}
	//err = editionRepo.SaveOne(&edition)
	//if err != nil {
	//	t.Errorf("could not save edition for saveEditionNew test - %v", err.Error())
	//}
	//
	//editionRepo.FindAll()

	//bookUuid := uuid.New().String()
	//record1Uuid := uuid.New().String()
	//record2Uuid := uuid.New().String()
	//page1Uuid := uuid.New().String()
	//page2Uuid := uuid.New().String()
	//
	book := Book{
		Id:           "testbookhash1",
		Name:         "testBook1",
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}

	record1 := Record{
		Id:               "testrecordhash1",
		FilePointer:      "whoCares",
		Name:             "testrecord1",
		Extension:        ".md",
		DateFileModified: time.Now(),
		DateCreated:      time.Now(),
		DateModified:     time.Now(),
	}

	record2 := Record{
		Id:               "testrecordhash2",
		FilePointer:      "whoCares2",
		Name:             "testrecord2",
		Extension:        ".md",
		DateFileModified: time.Now(),
		DateCreated:      time.Now(),
		DateModified:     time.Now(),
	}

	page1 := Page{
		Id:           "testpagehash1",
		RecordId:     record1.Id,
		Record:       &record1,
		EditionId:    "testEdition1",
		Edition:      nil,
		DateCreated:  time.Time{},
		DateModified: time.Time{},
	}

	page2 := Page{
		Id:           "testpagehash2",
		RecordId:     record2.Id,
		Record:       &record2,
		EditionId:    "testEdition1",
		Edition:      nil, //fixme this is sussy
		DateCreated:  time.Time{},
		DateModified: time.Time{},
	}

	edition := Edition{
		Id:            "testEdition1",
		EditionNumber: 0,
		BookId:        book.Id,
		Book:          &book,
		Pages:         &[]Page{page1, page2},
		DateCreated:   time.Now(),
		DateModified:  time.Now(),
	}

	t.Run("happypath", func(t *testing.T) {
		err := editionRepo.SaveOne(&edition)
		if err != nil {
			t.Errorf("couldnt save edition - %v", err.Error())
		}

		res, err := editionRepo.FindOne(edition.Id)
		if err != nil {
			t.Errorf("couldnt load the edition that was saved - %v", err.Error())
		}
		t.Logf("%#v", res)
	})
}

func deleteTestEditionEntries(db *gorm.DB) {
	db.Where("uuid like ?", "test%").Delete(&Edition{})
}

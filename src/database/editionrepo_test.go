package database

import (
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestEditionRepo_SaveOne(t *testing.T) {
	var db = Connect("../../out/library.db")

	var editionRepo = NewEditionRepo(db)
	var bookRepo = NewBookRepo(db)

	deleteTestEditionEntries(db)
	deleteTestBookEntries(db)
	defer deleteTestEditionEntries(db)
	defer deleteTestBookEntries(db)

	t.Run("saveEditionNew", func(t *testing.T) {
		var book = Book{
			Uuid:         "testBookuuid",
			Name:         "testBook",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		}
		err := bookRepo.SaveOne(&book)
		if err != nil {
			t.Errorf("faled to prep saveEditionNew test, could not save book - %v", err.Error())
		}

		var edition = Edition{
			Uuid:          "testEditionUuid",
			EditionNumber: 0,
			BookUuid:      "testBookuuid",
			Book:          nil,
			DateCreated:   time.Now(),
			DateModified:  time.Now(),
		}
		err = editionRepo.SaveOne(&edition)
		if err != nil {
			t.Errorf("could not save edition for saveEditionNew test - %v", err.Error())
		}
	})

	t.Run("saveEditionInvalidBookUuid", func(t *testing.T) {
		var edition = Edition{
			Uuid:          "testEditionUuidFailTarget",
			EditionNumber: 0,
			BookUuid:      "fakeuuid",
			Book:          nil,
			DateCreated:   time.Now(),
			DateModified:  time.Now(),
		}

		err := editionRepo.SaveOne(&edition)
		if err == nil {
			t.Errorf("completed edition save when should not be possible due to fkey: %v", edition.Book.Uuid)
		}
	})
}

func TestEditionRepo_FindOne(t *testing.T) {
	var db = Connect("../../out/library.db")

	var editionRepo = NewEditionRepo(db)
	var bookRepo = NewBookRepo(db)

	deleteTestEditionEntries(db)
	deleteTestBookEntries(db)
	//defer deleteTestBookEntries(db)
	//defer deleteTestEditionEntries(db)

	var book = Book{
		Uuid:         "testBookuuid",
		Name:         "testBook",
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}
	err := bookRepo.SaveOne(&book)
	if err != nil {
		t.Errorf("faled to prep saveEditionNew test, could not save book - %v", err.Error())
	}

	var edition = Edition{
		Uuid:          "testEditionUuid",
		EditionNumber: 0,
		BookUuid:      "testBookuuid",
		Book:          nil,
		DateCreated:   time.Now(),
		DateModified:  time.Now(),
	}
	err = editionRepo.SaveOne(&edition)
	if err != nil {
		t.Errorf("could not save edition for saveEditionNew test - %v", err.Error())
	}

	editionRepo.FindAll()
	t.Run("happypath", func(t *testing.T) {
		res, err := editionRepo.FindOne(edition.Uuid)
		if err != nil {
			t.Errorf("failed to fetch edition from db %v - %v", edition.Uuid, err.Error())
		}

		t.Logf("%#v", res)
	})
}

func deleteTestEditionEntries(db *gorm.DB) {
	db.Where("uuid like ?", "test%").Delete(&Edition{})
}

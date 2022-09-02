package service

import (
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/database/repository"
	"testing"
)

func TestEditionManagerService_CreateNewEditionInNamedBook_MANUALVERIFY(t *testing.T) {
	db := database.Connect("../../out/library_integration_test.db")
	repoManager := repository.NewRepoManager(db)
	fileService := NewFileService("../../out/filebin")
	editionService := NewEditionManagerService(repoManager, fileService)

	repository.WipeTestDatabase(db)
	//defer repository.WipeTestDatabase(db)

	book, err := editionService.CreateNewBook("testBook1")
	if err != nil {
		t.Fatalf("book create fail")
	}
	edition, err := editionService.CreateNewEdition(book.Name, "edition1", "../")
	if err != nil {
		t.Fatalf("edition create fail - %v", err.Error())
	}
	t.Logf("%#v", edition)
}

func TestEditionManagerService_CreateNewEdition(t *testing.T) {
	db := database.Connect("../../out/library_integration_test.db")
	repoManager := repository.NewRepoManager(db)
	fileService := NewFileService("../../out/filebin")
	editionService := NewEditionManagerService(repoManager, fileService)

	repository.WipeTestDatabase(db)
	//defer repository.WipeTestDatabase(db)

	book, err := editionService.CreateNewBook("testBook1")
	if err != nil {
		t.Fatalf("book create fail")
	}
	edition, err := editionService.CreateNewEdition(book.Name, "edition1", "../")
	if err != nil {
		t.Fatalf("edition create fail - %v", err.Error())
	}
	t.Logf("%#v", edition)

	_, err = editionService.CreateNewEdition("testBook1", "edition2", "../../schema")
	if err != nil {
		t.Fatalf("edition create - %v", err.Error())
	}
}

//func TestEditionManagerService_DownloadEdition(t *testing.T) {
//	//dependencies
//	db := database.Connect("../../out/library_integration_test.db")
//	repoManager := repository.NewRepoManager(db)
//	fileService := NewFileService("../../out/filebin")
//	editionService := NewEditionManagerService(repoManager, fileService)
//
//	//setup
//	repository.WipeTestDatabase(db)
//	//defer repository.WipeTestDatabase(db)
//
//	_, err := editionService.CreateNewBook("testBook1")
//	if err != nil {
//		t.Fatalf("book create fail")
//	}
//	edition, err := editionService.CreateNewEdition("testBook1", "edition0", "../")
//	if err != nil {
//		t.Fatalf("edition create fail - %v", err.Error())
//	}
//	t.Logf("%#v", edition)
//
//	_, err = editionService.CreateNewEdition("testBook1", "edition1", "../../schema")
//	if err != nil {
//		t.Fatalf("edition create - %v", err.Error())
//	}
//
//	_, err = editionService.CreateNewEdition("testBook1", "edition1", "C:\\Users\\alexl\\Downloads\\mangas\\Reborn as a Vending Machine, I Now Wander the Dungeon")
//	if err != nil {
//		t.Fatalf("edition create - %v", err.Error())
//	}
//
//	//exec
//	err = editionService.DownloadEdition("testBook1", 2, "../../out/recovertarget")
//	if err != nil {
//		t.Fatalf("failed op - %v", err.Error())
//	}
//}

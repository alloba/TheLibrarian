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

	//repository.WipeTestDatabase(db)
	//defer repository.WipeTestDatabase(db)

	edition, err := editionService.CreateNewEditionInNamedBook("testBook", "", "../")
	if err != nil {
		t.Fatalf("Could not form edition - %v", err.Error())
	}
	t.Logf("%#v", edition)
}

package main

import (
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/webserver"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

func main() {
	log.Println("Initializing the Librarian")

	var db = database.Connect("../schema/library.db")
	defer db.Close()

	var recordRepo = database.NewRecordRepo(db)
	testRepoOperation(recordRepo)

	log.Println("Terminating the Librarian")
}

func testRepoOperation(repo database.RecordRepo) {
	//var allRecords, err = repo.FindAll()
	//if err != nil {
	//	log.Fatalf("Couldnt do the thing: %v", err)
	//}
	//

	var record = database.Record{
		Hash:             "123",
		FilePointer:      "somelocation",
		Name:             "tst",
		Extension:        "aaa",
		DateFileModified: time.Now(),
		DateCreated:      time.Now(),
		DateModified:     time.Now(),
	}
	err := repo.SaveOne(&record)
	if err != nil {
		log.Fatalf("failed to save to the database. %v", err)
	}
	log.Println("testing repo operations has concluded")
}

func launchWebserver() {
	log.Println("Launching webserver")
	var server = webserver.New()

	server.ServerPort = ":8080"

	RegisterControllerEndpoints(&server)
	RegisterPreHooks(&server)
	RegisterPostHooks(&server)
	webserver.Start(&server)

	log.Println("Webserver running on port " + server.ServerPort)
}

package main

import (
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/webserver"
	"log"
)

func main() {
	log.Println("Initializing the Librarian")

	var db = database.Connect("../out/library.db")
	var recordRepo = database.NewRecordRepo(db)
	testRepoOperation(recordRepo)

	log.Println("Terminating the Librarian")
}

func testRepoOperation(repo *database.RecordRepo) {
	record, err := repo.FindByHash("123")
	if err != nil {
		log.Fatalf("couldnt find record: %v", err.Error())
	}
	log.Printf("Found record: %#v", record)
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

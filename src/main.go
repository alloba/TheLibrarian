package main

import (
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/logging"
	"github.com/alloba/TheLibrarian/webserver"
	"log"
)

//TODO i'd like a consolidated logTrace function. should probably move to it's own package...

func main() {
	log.Println("Initializing the Librarian")

	var db = database.Connect("../out/library_integration_test.db")
	//var recordRepo = database.NewRecordRepo(db)
	//testRepoOperation(recordRepo)
	coordinator := NewActionCoordinator(db, "../out/filebin/")
	err := coordinator.SubmitNewEdition("testBook1", "", "./")
	if err != nil {
		panic(logging.LogTrace(err))
	}

	log.Println("Terminating the Librarian")
}

func launchWebserver() {
	log.Println("Launching webserver")
	var server = webserver.New()

	server.ServerPort = ":8080"

	webserver.RegisterControllerEndpoints(&server)
	webserver.RegisterPreHooks(&server)
	webserver.RegisterPostHooks(&server)
	webserver.Start(&server)

	log.Println("Webserver running on port " + server.ServerPort)
}

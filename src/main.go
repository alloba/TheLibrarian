package main

import (
	"github.com/alloba/TheLibrarian/webserver"
	"log"
)

func main() {
	log.Println("Initializing the Librarian")

	//var db = database.Connect("../out/library.db")
	//var recordRepo = database.NewRecordRepo(db)
	//testRepoOperation(recordRepo)

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

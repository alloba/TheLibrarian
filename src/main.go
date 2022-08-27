package main

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/webserver"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	log.Println("Initializing the Librarian")

	var db = database.Connect("../schema/library.db")
	defer db.Close()

	var recordRepo = database.NewRecordRepo(db)
	var allRecords, err = recordRepo.FindAll()
	if err != nil {
		log.Fatalf("Couldnt do the thing: %v", err)
		//panic("couldnt do the thing")
	}
	fmt.Printf("%#v\n", allRecords)

	log.Println("Terminating the Librarian")
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

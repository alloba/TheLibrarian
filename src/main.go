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

	runTestInstance()

	log.Println("Terminating the Librarian")
}

func runTestInstance() {
	var db = database.Connect("C:\\Users\\alexl\\projects\\TheLibrarian\\out\\library_integration_test.db")
	coordinator := NewActionCoordinator(db, "C:\\Users\\alexl\\projects\\TheLibrarian\\out\\filebin\\")

	//err := coordinator.SubmitNewEdition("books", "", "~/books")
	//if err != nil {
	//	panic(logging.LogTrace(err))
	//}

	err := coordinator.DownloadEdition("books", 10, "../out/recovertarget")
	if err != nil {
		panic(logging.LogTrace(err))
	}
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

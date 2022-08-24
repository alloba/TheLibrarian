package main

import (
	"github.com/alloba/TheLibrarian/webserver"
	"log"
)

func main() {
	log.Println("Initializing the Librarian")

	var server = webserver.New()
	RegisterControllerEndpoints(&server)
	RegisterPreHooks(&server)
	RegisterPostHooks(&server)
	webserver.Start(&server)

	log.Println("Terminating the Librarian")
}

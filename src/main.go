package main

import (
	"database/sql"
	"github.com/alloba/TheLibrarian/webserver"
	"log"
)

type Env struct {
	db *sql.DB
}

func main() {
	log.Println("Initializing the Librarian")

	var server = webserver.New()
	RegisterControllerEndpoints(&server)
	RegisterPreHooks(&server)
	RegisterPostHooks(&server)
	webserver.Start(&server)

	log.Println("Terminating the Librarian")
}

//TODO
//	get this going -- https://github.com/mattn/go-sqlite3
//	create notes around running a local instance of sqlite
//	write down actual schema file for the library
//	best effort first attempt at writing to database
//  test functions for everything i'm writing
//  ... Everything Else.

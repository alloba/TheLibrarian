package main

import (
	"flag"
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/logging"
	"github.com/alloba/TheLibrarian/webserver"
	"log"
	"os"
	"syscall"
)

func main() {
	log.Println("Initializing the Librarian")

	flag.Usage = printHelpMessage
	submitFlag := flag.Bool("submit", false, "trigger submitting new edition")
	checkoutFlag := flag.Bool("checkout", false, "trigger downloading existing edition")
	folder := flag.String("folder", "", "target folder to operate on")
	book := flag.String("book", "", "target book to push an edition to")
	edition := flag.Int("edition", -1, "which edition to pull")
	libraryPath := flag.String("library", "", "path to library folder")
	flag.Parse()

	if !*submitFlag && !*checkoutFlag {
		printHelpMessage()
		syscall.Exit(1)
	}
	if *submitFlag && *checkoutFlag {
		printHelpMessage()
		syscall.Exit(1)
	}
	if *libraryPath == "" {
		printHelpMessage()
		syscall.Exit(1)
	}

	db := database.Connect(*libraryPath + string(os.PathSeparator) + "library.db")
	coordinator := NewActionCoordinator(db, *libraryPath)

	var err error
	if *submitFlag {
		if *book == "" || *folder == "" {
			printHelpMessage()
			syscall.Exit(1)
		}
		err = coordinator.SubmitNewEdition(*book, "", *folder)
	} else if *checkoutFlag {
		if *folder == "" || *book == "" {
			printHelpMessage()
			syscall.Exit(1)
		}
		if *edition == -1 {
			err = coordinator.DownloadNewestEdition(*book, *folder)
		} else {
			err = coordinator.DownloadEdition(*book, *edition, *folder)
		}
	}
	if err != nil {
		panic(logging.LogTrace(err))
	}
	log.Println("Terminating the Librarian")
}

func printHelpMessage() {
	fmt.Printf(
		"Librarian -- Help\n" +
			" This program takes a number of flags in to determine the specific action to take.\n" +
			" Primarily this means either providing a 'submit' or a 'checkout' flag (but not both).\n" +
			" --submit    : Add a new folder to the archive.\n" +
			"   --folder  : Specify which folder on the current machine should be submitted.\n" +
			"   --book    : Which book to save the folder underneath. It will be entered as a new edition for the named book.\n" +
			" --checkout  : Download an edition from the library. \n" +
			"   --folder  : The download destination for the edition. A subfolder will be placed in the provided directory. \n" +
			"   --book    : Which book to pull from." +
			"   --edition : (Optional) Which edition of the book to download. If not provided, will grab the newest edition." +
			"\n")
	//syscall.Exit(1)
}

/**
CLI Options:
librarian --submit --folder --book nameOfBook
librarian --checkout --book nameOfBook --edition 0 --folder ../

two primary options, with subsequent desired flags.
so could say if option invoke function with parameters. not bad.
*/

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

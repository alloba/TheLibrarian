package main

import (
	"fmt"
	"github.com/alloba/TheLibrarian/webserver"
	"net/http"
)

// RegisterControllerEndpoints handles adding all web endpoints to the server.
// This is hard-coded configured, and deals with private functions.
func RegisterControllerEndpoints(server *webserver.ServerManager) {
	for _, val := range testController {
		webserver.RegisterEndpoint(server, val)
	}
}

var testController = []webserver.Endpoint{
	{"/test/hello", webserver.GET, testControllerHelloWorld()},
	{"/test/counter", webserver.GET, testControllerGlobalCounter()},
}

func testControllerHelloWorld() webserver.RequestFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world!")
	}
}

var counter = 0

func testControllerGlobalCounter() webserver.RequestFunction {
	return func(writer http.ResponseWriter, request *http.Request) {
		counter += 1
		fmt.Fprint(writer, "Counted: ", counter)
	}
}

package webserver

import (
	"fmt"
	"net/http"
)

// RegisterControllerEndpoints handles adding all web endpoints to the server.
// This is hard-coded configured, and deals with private functions.
func RegisterControllerEndpoints(server *ServerManager) {
	for _, val := range testController {
		RegisterEndpoint(server, val)
	}
}

var testController = []Endpoint{
	{"/test/hello", GET, testControllerHelloWorld()},
	{"/test/counter", GET, testControllerGlobalCounter()},
}

func testControllerHelloWorld() RequestFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world!")
	}
}

var counter = 0

func testControllerGlobalCounter() RequestFunction {
	return func(writer http.ResponseWriter, request *http.Request) {
		counter += 1
		fmt.Fprint(writer, "Counted: ", counter)
	}
}

package webserver

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRegisterEndpoint(t *testing.T) {
	var servermanager = New()
	RegisterEndpoint(&servermanager, "testEndpoint", GET, func(writer http.ResponseWriter, request *http.Request) {})
	RegisterEndpoint(&servermanager, "testEndpoint", POST, func(writer http.ResponseWriter, request *http.Request) {})
	RegisterEndpoint(&servermanager, "testEndpoint2", GET, func(writer http.ResponseWriter, request *http.Request) {})
	RegisterEndpoint(&servermanager, "testEndpoint2", POST, func(writer http.ResponseWriter, request *http.Request) {})
	RegisterEndpoint(&servermanager, "testEndpoint3", GET, func(writer http.ResponseWriter, request *http.Request) {})
	RegisterEndpoint(&servermanager, "testEndpoint4", POST, func(writer http.ResponseWriter, request *http.Request) {})

	registerEndpointHandlers(&servermanager)
	fmt.Printf("%#v", servermanager)
}

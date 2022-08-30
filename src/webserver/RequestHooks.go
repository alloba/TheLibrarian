package webserver

import (
	"log"
	"net/http"
)

func logEnteringPath() RequestFunction {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Entering - %v %v", request.Method, request.RequestURI)
	}
}

func logLeavingPath() RequestFunction {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Exiting - %v %v", request.Method, request.RequestURI)
	}
}

func RegisterPreHooks(manager *ServerManager) {
	RegisterPreHandle(manager, logEnteringPath())
}

func RegisterPostHooks(manager *ServerManager) {
	RegisterPostHandle(manager, logLeavingPath())
}

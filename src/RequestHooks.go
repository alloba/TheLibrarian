package main

import (
	"github.com/alloba/TheLibrarian/webserver"
	"log"
	"net/http"
)

func logEnteringPath() webserver.RequestFunction {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Entering - %v %v", request.Method, request.RequestURI)
	}
}

func logLeavingPath() webserver.RequestFunction {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Exiting - %v %v", request.Method, request.RequestURI)
	}
}

func RegisterPreHooks(manager *webserver.ServerManager) {
	webserver.RegisterPreHandle(manager, logEnteringPath())
}

func RegisterPostHooks(manager *webserver.ServerManager) {
	webserver.RegisterPostHandle(manager, logLeavingPath())
}

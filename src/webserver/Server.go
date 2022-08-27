// This is a layer on top of net/http which adds a few conveniences.
// - HTTP Verb level function declaration
// - Pre- and Post-request handlers
// - Sensible defaults for configuration (currently only the port)

package webserver

import (
	"log"
	"net/http"
)

// HttpMethod Explicit type matching for supported HTTP methods, to prevent typos.
// Anything with this type is going to be directly matched to a Method constant from net/http
type HttpMethod string

// RequestFunction Any function that needs to match the standard net/http request function interface.
// This type is appropriate for both hooks and the endpoint operation itself.
type RequestFunction func(http.ResponseWriter, *http.Request)

// HTTP methods allowed for endpoints.
// These map directly to methods defined in net/http and are mainly for convenience.
const (
	GET    HttpMethod = http.MethodGet
	POST   HttpMethod = http.MethodPost
	PUT    HttpMethod = http.MethodPut
	DELETE HttpMethod = http.MethodDelete
)

// Endpoint A container for endpoint information.
// This struct provides enough information to merge multiple request functions onto a single endpoint (by Path + Method)
type Endpoint struct {
	Path             string
	Method           HttpMethod
	EndpointFunction RequestFunction
}

// ServerManager manages all information required to configure and run a webserver.
// This includes all relevant configuration as well as any registered hooks and endpoint functions.
type ServerManager struct {
	ServerPort      string                // ServerPort which port the server will run on.
	preHandleChain  []RequestFunction     // preHandleChain a list of functions that are run before any endpoint request
	postHandleChain []RequestFunction     // postHandleChain a list of functions that are run after any endpoint request
	endpoints       map[string][]Endpoint //endpoints all registered endpoints, organized by Path string
}

// New will instantiate a new instance of ServerManager
func New() ServerManager {
	return ServerManager{":8080", make([]RequestFunction, 0), make([]RequestFunction, 0), make(map[string][]Endpoint, 0)}
}

// Start triggers the webserver.
// This merges all pre- and post-request functions with all HTTP Method functions, and invokes http.ListenAndServe
func Start(manager *ServerManager) {
	registerEndpointHandlers(manager)
	log.Println("Starting server on port " + manager.ServerPort)
	if err := http.ListenAndServe(manager.ServerPort, nil); err != nil {
		log.Fatal(err)
	}
}

// RegisterPreHandle adds a new function to the chain .
func RegisterPreHandle(manager *ServerManager, hook func(http.ResponseWriter, *http.Request)) {
	manager.preHandleChain = append(manager.preHandleChain, hook)
}

// RegisterPostHandle adds a new function to the chain.
func RegisterPostHandle(manager *ServerManager, hook func(w http.ResponseWriter, r *http.Request)) {
	manager.postHandleChain = append(manager.postHandleChain, hook)
}

// RegisterEndpoint adds a new endpoint object to the web server.
// If the same path + method already exists in the server, the application will fail.
func RegisterEndpoint(manager *ServerManager, marker Endpoint) {
	if val, ok := manager.endpoints[marker.Path]; ok {
		for _, endpointObj := range val {
			if endpointObj.Method == marker.Method {
				log.Fatalf("Endpoint Method has already been registered: %v - %v", marker.Path, string(marker.Method))
			}
		}
		manager.endpoints[marker.Path] = append(manager.endpoints[marker.Path], marker)
	} else {
		manager.endpoints[marker.Path] = make([]Endpoint, 0)
		manager.endpoints[marker.Path] = append(manager.endpoints[marker.Path], marker)
	}
	//TODO: ensuring a leading slash is in the path would be nice (doesnt resolve otherwise)
}

// SetRunningPort assign the running port for the webserver
func SetRunningPort(manager *ServerManager, port string) {
	manager.ServerPort = port
}

// createEndpointFunction merges all hooks and all functions for an endpoint together.
// It also contains the logic that associates the HTTP Method with the particular function to execute
// We're on the honor system that the Endpoint list actually is all the same endpoint Path.
func createEndpointFunction(preHook []RequestFunction, postHook []RequestFunction, endpointMethodFunctions []Endpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, f := range preHook {
			f(w, r)
		}

		// checking the Method might be nicer with a map instead of a slice. future improvement maybe.
		for _, providedFunction := range endpointMethodFunctions {
			if string(providedFunction.Method) == r.Method {
				providedFunction.EndpointFunction(w, r)
				break
			}
		}

		for _, f := range postHook {
			f(w, r)
		}
	}
}

// registerEndpointHandlers coordinates createEndpointFunction.
// I could probably merge the two together.
func registerEndpointHandlers(manager *ServerManager) {
	for key, val := range manager.endpoints {
		log.Printf("Registering endpoints for: %v", key)
		var operatorFunction = createEndpointFunction(manager.preHandleChain, manager.postHandleChain, val)
		http.HandleFunc(key, operatorFunction)
	}
}

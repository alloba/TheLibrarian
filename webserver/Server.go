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
// Anything with this type is going to be directly matched to a method constant from net/http
type HttpMethod string

// requestFunction Any function that needs to match the standard net/http request function interface.
// This type is appropriate for both hooks and the endpoint operation itself.
type requestFunction func(http.ResponseWriter, *http.Request)

// HTTP methods allowed for endpoints.
// These map directly to methods defined in net/http and are mainly for convenience.
const (
	GET    HttpMethod = http.MethodGet
	POST   HttpMethod = http.MethodPost
	PUT    HttpMethod = http.MethodPut
	DELETE HttpMethod = http.MethodDelete
)

// endpointMarker A container for endpoint information.
// This struct provides enough information to merge multiple request functions onto a single endpoint (by path + method)
type endpointMarker struct {
	path             string
	method           HttpMethod
	endpointFunction requestFunction
}

// ServerManager manages all information required to configure and run a webserver.
// This includes all relevant configuration as well as any registered hooks and endpoint functions.
type ServerManager struct {
	serverPort      string                      // serverPort which port the server will run on.
	preHandleChain  []requestFunction           // preHandleChain a list of functions that are run before any endpoint request
	postHandleChain []requestFunction           // postHandleChain a list of functions that are run after any endpoint request
	endpoints       map[string][]endpointMarker //endpoints all registered endpoints, organized by path string
}

// New will instantiate a new instance of ServerManager
func New() ServerManager {
	return ServerManager{":8080", make([]requestFunction, 0), make([]requestFunction, 0), make(map[string][]endpointMarker, 0)}
}

// RegisterPreHandle adds a new function to the chain .
func RegisterPreHandle(manager *ServerManager, hook func(http.ResponseWriter, *http.Request)) {
	manager.preHandleChain = append(manager.preHandleChain, hook)
}

// RegisterPostHandle adds a new function to the chain.
func RegisterPostHandle(manager *ServerManager, hook func(w http.ResponseWriter, r *http.Request)) {
	manager.postHandleChain = append(manager.postHandleChain, hook)
}

// RegisterEndpoint adds a new endpoint to be managed by the server.
func RegisterEndpoint(manager *ServerManager, endpoint string, method HttpMethod, endpointFunction requestFunction) {
	if val, ok := manager.endpoints[endpoint]; ok {
		//if the method exists throw error
		for _, endpointObj := range val {
			if endpointObj.method == method {
				log.Fatalf("Endpoint method has already been registered: %v - %v", endpoint, method)
			}
		}
		//else add to list
		manager.endpoints[endpoint] = append(val, endpointMarker{endpoint, method, endpointFunction})
	} else {
		manager.endpoints[endpoint] = make([]endpointMarker, 0)
		manager.endpoints[endpoint] = append(manager.endpoints[endpoint], endpointMarker{endpoint, method, endpointFunction})
	}
}

// SetRunningPort assign the running port for the webserver
func SetRunningPort(manager *ServerManager, port string) {
	manager.serverPort = port
}

// createEndpointFunction merges all hooks and all functions for an endpoint together.
// It also contains the logic that associates the HTTP method with the particular function to execute
// We're on the honor system that the endpointMarker list actually is all the same endpoint path.
func createEndpointFunction(preHook []requestFunction, postHook []requestFunction, endpointMethodFunctions []endpointMarker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, f := range preHook {
			f(w, r)
		}

		// checking the method might be nicer with a map instead of a slice. future improvement maybe.
		for _, providedFunction := range endpointMethodFunctions {
			if string(providedFunction.method) == r.Method {
				providedFunction.endpointFunction(w, r)
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

// Start triggers the webserver.
// This merges all pre- and post-request functions with all HTTP method functions, and invokes http.ListenAndServe
func Start(manager *ServerManager) {
	registerEndpointHandlers(manager)
	log.Println("Starting server on port " + manager.serverPort)
	if err := http.ListenAndServe(manager.serverPort, nil); err != nil {
		log.Fatal(err)
	}
}

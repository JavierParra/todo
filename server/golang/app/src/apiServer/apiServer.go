// Package apiServer provides a wrapper around the http package to streamline
// the creation of JSON Restful api's.
package apiServer

import (
	"net/http"
	"fmt"
	"log"
	"reflect"
)

// RouteMethods serve as flags to indicate which methods a particular route
// shall respond to.
type RouteMethods struct {
	GET      bool
	POST     bool
	PUT      bool
	PATCH    bool
	DELETE   bool
	OPTIONS  bool
	HEAD     bool
	CONNECT  bool
	TRACE    bool
}

func NewFullMethods () *RouteMethods {
	return &RouteMethods{true, true, true, true, true, true, true, true, true}
}

type HandlerFunc func (server *Server, response *Response, r *Request)

type Server struct {
	registry map[string] map[string] HandlerFunc
	matcher  Matcher
}

func (server *Server) routerHandlerFactory (path string) http.HandlerFunc {
	return func (writer http.ResponseWriter, req *http.Request) {
		method := req.Method
		handler := server.registry[path][method]
		response := Response{ writer: writer }
		request := Request{ request: req }

		if handler == nil {
			response.SendWithStatus(&ApiError{
				"BAD_METHOD",
				fmt.Sprintf("Method %s not allowed", method),
				struct{ Method string `json:"method"` }{ method },
			}, http.StatusMethodNotAllowed)
			return
		}

		handler(server, &response, &request)
	}
}

func (server *Server) registerRoute (path string, methods *RouteMethods, handler HandlerFunc) {
	patternId := server.matcher.RegisterPattern(path)
	// fmt.Println(patternId, path)
	if server.registry[patternId] == nil {
		server.registry[patternId] = make(map[string] HandlerFunc, 0)
		// Registers function handler for path in go's http server.
		http.HandleFunc(path, server.routerHandlerFactory(path))
	}

	routes := server.registry[patternId]

	reflection := reflect.Indirect(reflect.ValueOf(methods))
	typeOfS := reflection.Type()

	for i := 0; i< typeOfS.NumField(); i++ {
		key := typeOfS.Field(i).Name
		val := reflection.Field(i).Interface()
		if val == false {
			continue
		}

		if routes[key] != nil {
			log.Fatalf("ERROR: Registering duplicated route for path '%s' and method '%s'", path, key)
		}
		routes[key] = handler
	}
}

// Route registers a function that shall respond to specific methods in a
// specific path. The path is parsed in the same way as net/http.
func (server *Server) Route (path string, methods *RouteMethods, handler HandlerFunc) bool {
	server.registerRoute(path, methods, handler)
	return true
}

// RouteAll registers a function that shall respondo to all methods in a
// specific path. The path is parsed in the same way as net/http.
func (server *Server) RouteAll (path string, handler HandlerFunc) bool {
	server.registerRoute(path, NewFullMethods(), handler)
	return true
}

func (server *Server) Serve (address string) {
	fmt.Println("Servng in", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func NewServer () (Server) {
	server := Server{
		registry: make(map[string] map[string] HandlerFunc, 0),
		matcher: Matcher{},
	}

	return server
}

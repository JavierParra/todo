package main

import (
	"fmt"
	"net/http"
	"time"

	"server/todo"
	"server/apiServer"
)

type ApiError = apiServer.ApiError
var collection = todo.GetStore()

func handler(server *apiServer.Server, response *apiServer.Response, request *apiServer.Request) {
	type List struct {
		Result []*todo.Todo `json:"todos"`
	}

	result := List{todo.Values(collection.Collection)}

	response.Send(result)
}

func create(server *apiServer.Server, response *apiServer.Response, request *apiServer.Request) {
	var todo todo.Todo
	err := request.ReadInto(&todo)

	if (request.Method() != http.MethodPost) {
		return
	}

	if err != nil {
		fmt.Println(err)
		response.SendWithStatus(&ApiError{Error: "VALIDATION_ERROR", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	todo.Created = time.Now().Unix()

	if todo.Name == "" {
		response.SendWithStatus(&ApiError{Error: "VALIDATION_ERROR", Message: "name is required"}, http.StatusBadRequest)
		return
	}

	idChan := make(chan string)
	errorChan := make(chan error)

	go collection.Add(&todo, idChan, errorChan)

	select {
		case <- idChan:
			response.SendWithStatus(&todo, http.StatusCreated)
		case err := <- errorChan:
			response.SendWithStatus(&ApiError{Error: "ERROR_RANDOM", Message: err.Error()}, 503)
	}
}

func get(server *apiServer.Server, response *apiServer.Response, request *apiServer.Request) {
	matches := request.GetMatches()
	id := matches["id"]
	todo := collection.Get(id)

	if todo == nil {
		response.SendWithStatus(&apiServer.ApiError{
			"NOT_FOUND",
			"The requested document was not found",
			struct{ Id string `json"id"`}{ id },
		}, http.StatusNotFound)
		return
	}

	response.Send(todo)
}


func deleteTodo(server *apiServer.Server, response *apiServer.Response, request *apiServer.Request) {
	matches := request.GetMatches()
	id := matches["id"]
	todo := collection.Get(id)

	if todo == nil {
		response.SendWithStatus(&apiServer.ApiError{
			"NOT_FOUND",
			"The requested document was not found",
			struct{ Id string `json"id"`}{ id },
		}, http.StatusNotFound)
		return
	}

	collection.Delete(id)

	response.Send(struct{
		Id      string `json"id"`
		Deleted bool   `json"deleted"`
	} { id, true })
}

func main() {
	address := "0.0.0.0:8000"
	server := apiServer.NewServer()
	server.Route("/todos", &apiServer.RouteMethods{POST: true}, create)
	server.Route("/todos", &apiServer.RouteMethods{GET: true}, handler)
	server.Route("/todos/:id([a-f0-9\\-]+-[a-f0-9\\-]+)", &apiServer.RouteMethods{GET: true}, get)
	server.Route("/todos/:id([a-f0-9\\-]+-[a-f0-9\\-]+)", &apiServer.RouteMethods{DELETE: true}, deleteTodo)
	server.Serve(address)
}

package main

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"

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
			response.Send(&todo)
		case err := <- errorChan:
			response.SendWithStatus(&ApiError{Error: "ERROR_RANDOM", Message: err.Error()}, 503)
	}
}

func get(server *apiServer.Server, response *apiServer.Response, request *apiServer.Request) {
	val, err := request.Read()

	if err != nil {
		response.SendInternalError(err)
		return
	}

	fmt.Println(val)
}

func sendTodo(writer http.ResponseWriter, todo *todo.Todo) {
	body, err := json.Marshal(todo)
	status := http.StatusOK

	if err != nil {
		fmt.Println("Unable to format json", err)
		status = http.StatusPartialContent
		body = []byte(`{}`)
	}

	writer.Header().Add("Content-type", "application/json")
	writer.WriteHeader(status)
	writer.Write(body)
}

func sendError(writer http.ResponseWriter, status int, error ApiError) {
	body, err := json.Marshal(error)

	if err != nil {
		fmt.Println("Unable to format json", err)
		body = []byte(`{}`)
	}

	writer.Header().Add("Content-type", "application/json")
	writer.WriteHeader(status)
	writer.Write(body)
}

func main() {
	address := "0.0.0.0:8000"
	server := apiServer.NewServer()
	server.Route("/todos", &apiServer.RouteMethods{POST: true}, create)
	server.Route("/todos", &apiServer.RouteMethods{GET: true}, handler)
	server.Route("/todos/:id", &apiServer.RouteMethods{GET: true}, handler)
	server.Route("/todos/:id([a-f0-9\\-]+-[a-f0-9\\-]+)/sortby/:sort([a-z])", &apiServer.RouteMethods{GET: true}, handler)
	server.Serve(address)
}

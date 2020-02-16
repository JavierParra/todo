package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/nu7hatch/gouuid"

	// "time"
	"encoding/json"
	"errors"
	"server/apiServer"
)

type Todo struct {
	Id        string `json:"id"`
	Complete  bool   `json:"complete"`
	Name      string `json:"name"`
	Created   int64  `json:"created"`
	Completed int64  `json:"completed"`
	Notes     string `json:"notes"`
}

type Store struct {
	Collection map[string]*Todo
}

func keys (col map[string]*Todo) []string {
	var i int = 0
	keys := make([]string, len(col))

	for k := range col {
		keys[i] = k
		i += 1
	}

	return keys
}

func values (col map[string]*Todo) []*Todo {
	var i int = 0
	keys := make([]*Todo, len(col))

	for _, v := range col {
		keys[i] = v
		i += 1
	}

	return keys
}

var collection = &Store{Collection: make(map[string]*Todo)}

func (store *Store) add (todo *Todo, idChan chan string, errorChan chan error) {
	uid, err := uuid.NewV4()
	id := uid.String()

	if err != nil {
		errorChan <- err
	}

	if rand.Int() % 23 == 0 {
		errorChan <- errors.New("Random error happened")
	}

	todo.Id = id

	store.Collection[id] = todo
	idChan <- id
}

func handler(server *apiServer.Server, response *apiServer.Response, request *apiServer.Request) {
	type List struct {
		Result []*Todo `json:"todos"`
	}

	result := List{values(collection.Collection)}

	response.Send(result)
}

func create(server *apiServer.Server, response *apiServer.Response, request *apiServer.Request) {
	var todo Todo
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

	go collection.add(&todo, idChan, errorChan)

	select {
		case <- idChan:
			response.Send(&todo)
		case err := <- errorChan:
			response.SendWithStatus(&ApiError{Error: "ERROR_RANDOM", Message: err.Error()}, 503)
	}
}

type ApiError = apiServer.ApiError

func sendTodo(writer http.ResponseWriter, todo *Todo) {
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
	server.Serve(address)
}

package main

import (
	"fmt"
	"io/ioutil"
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

func handler(writer http.ResponseWriter, request *http.Request) {
	type List struct {
		Result []*Todo `json:"todos"`
	}

	result := List{values(collection.Collection)}

	res, _ := json.Marshal(result)
	writer.Header().Add("Content-type", "application/json")
	writer.Write(res)
}

func create(writer http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodPost) {
		sendError(writer, http.StatusMethodNotAllowed, ApiError{"BAD_METHOD", "Method not allowd"})
		return
	}

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		fmt.Println(err)
		sendError(writer, http.StatusInternalServerError, ApiError{"INTERNAL", "Internal server error"})
		return
	}

	var todo Todo
	err = json.Unmarshal(body, &todo)

	if err != nil {
		fmt.Println(err)
		sendError(writer, http.StatusBadRequest, ApiError{"VALIDATION_ERROR", err.Error()})
		return
	}

	todo.Created = time.Now().Unix()

	if todo.Name == "" {
		sendError(writer, http.StatusBadRequest, ApiError{"VALIDATION_ERROR", "name is required"})
		return
	}

	idChan := make(chan string)
	errorChan := make(chan error)

	go collection.add(&todo, idChan, errorChan)

	select {
		case <- idChan:
			sendTodo(writer, &todo)
		case err := <- errorChan:
			sendError(writer, 503, ApiError{"ERROR_RANDOM", err.Error()})
	}
}

type ApiError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

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

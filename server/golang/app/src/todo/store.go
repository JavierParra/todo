package todo

import (
	"github.com/nu7hatch/gouuid"
	"errors"
	"math/rand"
)

type Store struct {
	Collection map[string]*Todo
}

func Keys (col map[string]*Todo) []string {
	var i int = 0
	keys := make([]string, len(col))

	for k := range col {
		keys[i] = k
		i += 1
	}

	return keys
}

func Values (col map[string]*Todo) []*Todo {
	var i int = 0
	keys := make([]*Todo, len(col))

	for _, v := range col {
		keys[i] = v
		i += 1
	}

	return keys
}

func (store *Store) Add (todo *Todo, idChan chan string, errorChan chan error) {
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

func (store *Store) Get (id string) *Todo {
	return store.Collection[id]
}

func GetStore() *Store {
	return &Store{Collection: make(map[string]*Todo)}
}

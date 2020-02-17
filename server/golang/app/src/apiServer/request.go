package apiServer

import (
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
)

// Request parses and validates a JSON request.
type Request struct {
	request *http.Request
	body    interface{}
	rawBody []byte
	bodyErr error
	matches map[string] string
}

// Read the body of a request once and returns it whenever it's called.
func (request *Request) Read () (interface{}, error) {
	if request.bodyErr != nil {
		return nil, request.bodyErr
	}

	if request.body != nil {
		return request.body, nil
	}

	body, err := ioutil.ReadAll(request.request.Body)

	if err != nil {
		fmt.Println(err)
		request.bodyErr = err
		return nil, err
	}

	request.rawBody = body

	err = json.Unmarshal(body, &request.body)

	return request.body, err
}

func (request *Request) ReadInto(v interface{}) error {
	_, err := request.Read()

	if err != nil {
		return err
	}

	return json.Unmarshal(request.rawBody, &v)
}

func (request *Request) GetMatches() map[string] string {
	return request.matches
}

func (request *Request) setMatches(matches map[string] string) {
	if request.matches != nil {
		log.Fatal("Trying to set matches more than once")
	}

	request.matches = matches
}

func (request *Request) Method () string {
	return request.request.Method
}

package apiServer

import (
	"net/http"
	"encoding/json"
	"fmt"
)

type Response struct {
	writer http.ResponseWriter
}

func (response *Response) Send(body interface{}) {
	response.SendWithStatus(body, http.StatusOK)
}

func (response *Response) SendWithStatus(body interface{}, status int) {
	resp, err := json.Marshal(body)
	writer := response.writer

	if err != nil {
		fmt.Println("Unable to format json", err)
		status = http.StatusPartialContent
		resp = []byte(`{}`)
	}

	fmt.Printf("Sending: %s \n", resp)

	writer.Header().Add("Content-type", "application/json")
	writer.WriteHeader(status)
	writer.Write(resp)
}

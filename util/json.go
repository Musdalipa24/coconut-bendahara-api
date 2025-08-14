package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result interface{}) (error) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	if err != nil {
		return err
	}

	return err
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	if err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		return
	}
}
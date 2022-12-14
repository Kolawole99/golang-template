package helper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// Response is used for a Static shape of the return response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

// EmptyObj is used when data does not what to be null on response object
type EmptyObj struct{}

const (
	Success string = "success"
	Failure string = "failure"
)

// BuildSuccessResponse is to inject data to dynamically handle success response
func BuildSuccessResponse(message string, data interface{}) Response {
	res := Response{
		Status:  Success,
		Message: message,
		Errors:  nil,
		Data:    data,
	}

	return res
}

// BuildErrorResponse is to inject data to dynamically handle error response
func BuildErrorResponse(message string, errors error) Response {
	splittedErrors := strings.Split(errors.Error(), "\n")

	res := Response{
		Status:  Failure,
		Message: message,
		Errors:  splittedErrors,
		Data:    nil,
	}

	return res
}

// WriteJSON takes a response status code and arbitrary data and writes a json response to the client
func WriteJSON(w http.ResponseWriter, status int, data Response, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// ReadJSON tries to read the body of a request and converts it into JSON
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

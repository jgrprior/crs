package crs

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// JSONError is used to construct simple JSON HTTP responses.
type JSONError struct {
	Status   int      `json:"status"`
	Messages []string `json:"messages"`
}

// JSONEncode marshals HTTP response messages into JSON.
func (e *JSONError) JSONEncode() (string, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(e)
	return buffer.String(), err
}

// JSONErrorResponse writes message and status to an http.ResponseWriter
func JSONErrorResponse(w http.ResponseWriter, msg string, stat int) {
	e := JSONError{stat, []string{msg}}
	s, _ := e.JSONEncode()
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(stat)
	io.WriteString(w, s)
}

// JSONErrorResponses writes message and status to an http.ResponseWriter
func JSONErrorResponses(w http.ResponseWriter, msg []string, stat int) {
	e := JSONError{stat, msg}
	s, _ := e.JSONEncode()
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(stat)
	io.WriteString(w, s)
}

// JSONResponse is used to construct simple JSON HTTP responses.
type JSONResponse struct {
	Status  int    `json:"status"`
	EntryID string `json:"entryId"`
}

// JSONEncode marshals HTTP response messages into JSON.
func (r *JSONResponse) JSONEncode() (string, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(r)
	return buffer.String(), err
}

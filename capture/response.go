package capture

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type JsonError struct {
	Status   int      `json:"status"`
	Messages []string `json:"messages"`
}

func (e *JsonError) JsonEncode() (string, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(e)
	return buffer.String(), err
}

// JsonErrorResponse writes message and status to an http.ResponseWriter
func JsonErrorResponse(w http.ResponseWriter, msg string, stat int) {
	e := JsonError{stat, []string{msg}}
	s, _ := e.JsonEncode()
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(stat)
	io.WriteString(w, s)
}

// JsonErrorResponse writes message and status to an http.ResponseWriter
func JsonErrorResponses(w http.ResponseWriter, msg []string, stat int) {
	e := JsonError{stat, msg}
	s, _ := e.JsonEncode()
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(stat)
	io.WriteString(w, s)
}

type JsonResponse struct {
	Status  int    `json:"status"`
	EntryId string `json:"entryId"`
}

func (r *JsonResponse) JsonEncode() (string, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(r)
	return buffer.String(), err
}

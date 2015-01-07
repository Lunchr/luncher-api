package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Handler func(http.ResponseWriter, *http.Request) *handlerError

type handlerError struct {
	Error   error
	Message string
	Code    int
}

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Print(e.Error)
		http.Error(w, e.Message, e.Code)
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) *handlerError {
	data, err := json.Marshal(v)
	if err != nil {
		return &handlerError{err, "", http.StatusInternalServerError}
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return nil
}

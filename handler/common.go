package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	. "github.com/Lunchr/luncher-api/router"
)

func writeJSON(w http.ResponseWriter, v interface{}) *HandlerError {
	data, err := json.Marshal(v)
	if err != nil {
		return &HandlerError{err, "", http.StatusInternalServerError}
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return nil
}

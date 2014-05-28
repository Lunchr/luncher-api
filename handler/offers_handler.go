package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func Offers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, Response{"response": "Hello, %s!"}.String(), html.EscapeString(r.URL.Path))
}

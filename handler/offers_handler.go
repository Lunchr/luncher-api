package handler

import (
	"fmt"
	"html"
	"net/http"
)

func Offers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

package main

import (
	"fmt"
	// "github.com/gorilla/mux"
	"html"
	"log"
	"net/http"
)

func main() {
	// r := mux.NewRouter()
	// r.HandleFunc("/offers", OffersHandler)
	// http.HandleFunc("/api", r)
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

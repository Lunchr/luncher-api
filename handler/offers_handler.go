package handler

import (
	"encoding/json"
	"fmt"
	"github.com/deiwin/praad-api/db"
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
	offers := []interface{}{}
	var offer interface{}
	iter := db.Offers.Find(struct{}{}).Iter()
	for {
		if iter.Next(&offer) {
			offers = append(offers, offer)
		} else {
			break
		}
	}
	fmt.Fprint(w, Response{"offers": offers})
}

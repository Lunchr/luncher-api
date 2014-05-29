package handler

import (
	"github.com/deiwin/praad-api/db"
	"github.com/deiwin/praad-api/db/model"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func Offers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var offers []*model.Offer
	var offer *model.Offer
	iter := db.Offers.Find(bson.M{}).Iter()
	for {
		if iter.Next(&offer) {
			offers = append(offers, offer)
		} else {
			break
		}
	}
	writeJson(w, offers)
}

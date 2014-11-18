package handler

import (
	"net/http"

	"log"

	"github.com/deiwin/praad-api/db"
	"github.com/deiwin/praad-api/db/model"
	"labix.org/v2/mgo/bson"
)

func Offers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var offers []*model.Offer
	err := db.Database.C(model.OfferCollectionName).Find(bson.M{}).All(&offers)
	if err != nil {
		log.Println(err)
	}
	writeJson(w, offers)
}

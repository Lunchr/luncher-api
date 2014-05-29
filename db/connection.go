package db

import (
	"labix.org/v2/mgo"
)

var (
	Offers  *mgo.Collection
	session *mgo.Session
)

func Connect() {
	var err error
	session, err = mgo.Dial("mongodb://localhost/test")
	if err != nil {
		panic(err)
	}
	db := session.DB("test")
	Offers = db.C("offers")
}

func Disconnect() {
	session.Close()
}

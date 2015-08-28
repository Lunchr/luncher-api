package model

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

const RegistrationAccessTokenCollectionName = "tags"

type (
	RegistrationAccessToken struct {
		ID    bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
		Token Token         `json:"token"         bson:"token"`
	}

	Token [16]byte
)

func NewToken() (Token, error) {
	var t [16]byte
	_, err := rand.Read(t[:])
	if err != nil {
		return Token{}, err
	}
	return Token(t), nil
}

func TokenFromString(s string) (Token, error) {
	var t Token
	tokenWithoutDashes := strings.Replace(s, "-", "", 4)
	b, err := hex.DecodeString(tokenWithoutDashes)
	if err != nil {
		return Token{}, err
	} else if len(b) != len(t) {
		return Token{}, errors.New("Token string is of unexpected length")
	}
	copy(t[:], b)
	return t, nil
}

func (t Token) String() string {
	return fmt.Sprintf("%X-%X-%X-%X-%X", t[0:4], t[4:6], t[6:8], t[8:10], t[10:])
}

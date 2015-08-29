package main

import (
	"fmt"
	"os"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
)

type RegistrationToken struct {
	Collection db.RegistrationAccessTokens
}

func (r RegistrationToken) CreateAndAdd() {
	token, err := model.NewRegistrationAccessToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if _, err = r.Collection.Insert(token); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Token %s successfully created and added!\n")
}

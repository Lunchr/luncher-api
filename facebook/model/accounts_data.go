package model

// Accounts is used to unmarshal the /me/accounts response
type Accounts struct {
	Data []AccountsData `json:"data"`
}

// AccountsData is the array element in the /me/accounts response
type AccountsData struct {
	ID          string `json:"id"`
	AccessToken string `json:"access_token"`
}

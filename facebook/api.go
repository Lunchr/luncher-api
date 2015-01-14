package facebook

import "net/http"

const (
	apiVersion = "v2.2"
)

type API interface {
	Connection(*http.Client) Connection
}

type api struct {
	graphURL string
}

func NewAPI() API {
	return api{
		graphURL: "https://graph.facebook.com/" + apiVersion,
	}
}

func (a api) Connection(client *http.Client) Connection {
	return connection{
		api:    a,
		Client: client,
	}
}

package facebook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/deiwin/luncher-api/facebook/model"
)

// Connection provides access to the Facebook API graph methods
type Connection interface {
	Me() (model.User, error)
}

type connection struct {
	*http.Client
	api API
}

func NewConnection(api API, client *http.Client) Connection {
	return connection{
		api:    api,
		Client: client,
	}
}

func (c connection) Me() (user model.User, err error) {
	resp, err := c.get("/me")
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &user)
	return
}

func (c connection) get(path string) (response []byte, err error) {
	resp, err := c.Get(c.api.graphURL + path)
	if err != nil {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		return
	}
	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	return
}

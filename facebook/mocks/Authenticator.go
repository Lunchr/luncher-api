package mocks

import "github.com/deiwin/facebook"
import "github.com/stretchr/testify/mock"

import "net/http"
import "golang.org/x/oauth2"

type Authenticator struct {
	mock.Mock
}

func (m *Authenticator) AuthURL(state string) string {
	ret := m.Called(state)

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Authenticator) Token(state string, r *http.Request) (*oauth2.Token, error) {
	ret := m.Called(state, r)

	var r0 *oauth2.Token
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*oauth2.Token)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Authenticator) APIConnection(tok *oauth2.Token) facebook.API {
	ret := m.Called(tok)

	r0 := ret.Get(0).(facebook.API)

	return r0
}
func (m *Authenticator) PageAccessToken(tok *oauth2.Token, pageID string) (string, error) {
	ret := m.Called(tok, pageID)

	r0 := ret.Get(0).(string)
	r1 := ret.Error(1)

	return r0, r1
}

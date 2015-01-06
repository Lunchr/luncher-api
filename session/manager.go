package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
)

const sessionCookieName = "luncher_session"

// Manager for the current session
type Manager interface {
	// GetOrInitSession returns the current session ID stored in the request
	// as a cookie or creates a new id and writes it into the response as a cookie
	GetOrInitSession(http.ResponseWriter, *http.Request) string
}

type manager struct{}

// NewManager returns an implementation of the Manager interface
func NewManager() Manager {
	return manager{}
}

func (mgr manager) GetOrInitSession(w http.ResponseWriter, r *http.Request) (session string) {
	if cookie, err := r.Cookie(sessionCookieName); err == nil && cookie.Value != "" {
		session = cookie.Value
	} else {
		session = createNewSession()
		cookie := &http.Cookie{
			Name:  sessionCookieName,
			Value: session,
			// no MaxAge because we want this to be a session cookie
		}
		http.SetCookie(w, cookie)
	}
	return
}

// createNewSession creates a new random session ID string
func createNewSession() string {
	// 16 bytes should be quite enough
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("We ran out of random!?")
	}
	return base64.URLEncoding.EncodeToString(b)
}

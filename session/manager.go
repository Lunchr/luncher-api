package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
)

const sessionCookieName = "luncher_session"

var ErrNotFound = errors.New("session manager: no session found")

// Manager for the current session
type Manager interface {
	// GetOrInit returns the current session ID stored in the request
	// as a cookie or creates a new id and writes it into the response as a cookie
	GetOrInit(http.ResponseWriter, *http.Request) string
	// Get returns the current session only if it exists. Returns an error otherwise
	Get(*http.Request) (string, error)
}

type manager struct{}

// NewManager returns an implementation of the Manager interface
func NewManager() Manager {
	return manager{}
}

func (m manager) Get(r *http.Request) (string, error) {
	if cookie, err := r.Cookie(sessionCookieName); err == http.ErrNoCookie {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	} else if cookie.Value == "" {
		return "", ErrNotFound
	} else {
		return cookie.Value, nil
	}
}

func (m manager) GetOrInit(w http.ResponseWriter, r *http.Request) string {
	if cookie, err := r.Cookie(sessionCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	session := createNewSession()
	cookie := &http.Cookie{
		Name:  sessionCookieName,
		Value: session,
		Path:  "/",
		// no MaxAge because we want this to be a session cookie
	}
	http.SetCookie(w, cookie)
	return session
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

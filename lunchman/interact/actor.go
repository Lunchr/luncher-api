// Package interact provides ways to interact with the user
package interact

import "io"

// An Actor provides methods to interact with the user
type Actor struct {
	rd io.Reader
}

// New creates a new Actor instance with the specified io.Reader
func NewActor(rd io.Reader) Actor {
	return Actor{rd}
}

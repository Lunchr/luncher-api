// Package interact provides ways to interact with the user
package interact

import "io"

type Interact struct {
	rd io.Reader
}

// New creates a new Interact instance with the specified io.Reader
func New(rd io.Reader) Interact {
	return Interact{rd}
}

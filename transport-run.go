// Copyright (c) 2017, Mitchell Cooper
package wikiclient

// RunTransport is a JSON transport over standard IO.
// it runs the wikiserver as a subprocess of quiki.

import (

)

type RunTransport struct {
	*jsonTransport
	wikiserverPath string
}

// create
func NewRunTransport(path string) (*RunTransport, error) {
	return &RunTransport{
		createJson(),
		path,
	}, nil
}

// connect
func (tr *RunTransport) Connect() error {
    return nil
}

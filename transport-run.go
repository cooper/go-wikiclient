// Copyright (c) 2017, Mitchell Cooper
package wikiclient

// RunTransport is a JSON transport over standard IO.
// it runs the wikiserver as a subprocess of quiki.

import (
	"bufio"
	"os/exec"
)

type RunTransport struct {
	*jsonTransport
	wikifierPath string
	configPath   string
	cmd          *exec.Cmd
}

// create
func NewRunTransport(wikifierPath, configPath string) (*RunTransport, error) {
	return &RunTransport{
		createJson(),
		wikifierPath,
		configPath,
		nil,
	}, nil
}

// connect
func (tr *RunTransport) Connect() error {

	// create readwriter for stdio
	var reader bufio.Reader
	var writer bufio.Writer
	tr.reader = &reader
	tr.writer = &writer
	rw := bufio.NewReadWriter(&reader, &writer)

	// create command
	cmd := exec.Command("wikiserver", "--std", tr.configPath)
	cmd.Dir = tr.wikifierPath
	cmd.Stdout = rw
	cmd.Stdin = rw
	tr.cmd = cmd

	// start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	tr.connected = true
	tr.startLoops()
	return nil
}

func (tr *RunTransport) startLoops() {
	// TODO: run continuously
	go tr.cmd.Wait()
	tr.jsonTransport.startLoops()
}

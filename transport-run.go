// Copyright (c) 2017, Mitchell Cooper
package wikiclient

// RunTransport is a JSON transport over standard IO.
// it runs the wikiserver as a subprocess of quiki.

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
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

	// try to find the absolute path of the config
	var cfg string
	var err error
	if cfg, err = filepath.Abs(tr.configPath); err != nil {
		cfg = tr.configPath
	}

	// create command
	cmd := exec.Command("./wikiserver", "--std", cfg)
	cmd.Dir = tr.wikifierPath
	tr.cmd = cmd

	// create stdio pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	tr.reader = bufio.NewReader(stdout)
	tr.writer = stdin

	// redirect stderr to ours
	cmd.Stderr = os.Stderr

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

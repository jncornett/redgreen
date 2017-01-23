package main

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/jncornett/redgreen"
	"github.com/jncornett/restful"
	"github.com/urfave/cli"
)

func doRun(c *cli.Context) error {
	if c.NArg() < 2 {
		if c.NArg() < 1 {
			return fmt.Errorf("missing key name")
		}
		return fmt.Errorf("missing command name")
	}
	args := c.Args()
	key := args[0]
	cmdName := args[1]
	cmdArgs := args[2:]
	lines, ok, cmdErr := execCommand(cmdName, cmdArgs, c.Bool("redir"))
	client := getClient(c.String("addr"))
	_, clientErr := client.Put(redgreen.Entry{
		ID:   restful.ID(key),
		OK:   ok,
		Data: lines,
	})
	if cmdErr == nil {
		return clientErr
	}
	return cmdErr
}

func execCommand(name string, args []string, redir bool) (lines []string, ok bool, err error) {
	cmd := exec.Command(name, args...)
	var out []byte
	if redir {
		out, err = cmd.CombinedOutput()
	} else {
		out, err = cmd.Output()
	}
	ok = err == nil
	var linesErr error
	lines, linesErr = readAllLines(bytes.NewReader(out))
	if err == nil && linesErr != nil {
		err = linesErr
	}
	return
}

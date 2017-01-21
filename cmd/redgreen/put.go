package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jncornett/redgreen"
	"github.com/urfave/cli"
)

func doPut(c *cli.Context) error {
	if c.NArg() < 2 {
		if c.NArg() < 1 {
			return fmt.Errorf("missing key") // TODO unclear
		}
		return fmt.Errorf("missing value") // TODO unclear
	}
	args := c.Args()
	key := args[0]
	ok, err := strconv.ParseBool(args[1])
	if err != nil {
		return err
	}
	lines, err := readAllLines(os.Stdin)
	if err != nil {
		return err
	}
	client := getClient(c.String("addr"))
	// TODO add strings from STDIN
	return client.Put(redgreen.Entry{
		Key:  key,
		OK:   ok,
		Data: lines,
	})
}

package main

import (
	"fmt"

	"github.com/jncornett/redgreen"
	"github.com/urfave/cli"
)

func doPop(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("missing key name")
	}
	key := c.Args()[0]
	client := getClient(c.String("addr"))
	entry, ok, err := client.Pop(key)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("key not found")
	}
	return printEntries([]redgreen.Entry{entry})
}

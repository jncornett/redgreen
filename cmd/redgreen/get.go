package main

import (
	"fmt"

	"github.com/jncornett/redgreen"
	"github.com/jncornett/restful"
	"github.com/urfave/cli"
)

func doGet(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("missing key name")
	}
	key := c.Args()[0]
	client := getClient(c.String("addr"))
	v, err := client.Get(restful.ID(key))
	if err != nil {
		return err
	}
	return printEntries([]redgreen.Entry{*v.(*redgreen.Entry)})
}

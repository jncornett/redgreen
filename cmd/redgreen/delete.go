package main

import (
	"fmt"

	"github.com/jncornett/restful"
	"github.com/urfave/cli"
)

func doDelete(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("missing key name")
	}
	key := c.Args()[0]
	client := getClient(c.String("addr"))
	err := client.Delete(restful.ID(key))
	if err != nil {
		return err
	}
	return nil
}

package main

import (
	"github.com/jncornett/redgreen"
	"github.com/urfave/cli"
)

func doGetAll(c *cli.Context) error {
	client := getClient(c.String("addr"))
	v, err := client.GetAll()
	if err != nil {
		return err
	}
	return printEntries(*v.(*[]redgreen.Entry))
}

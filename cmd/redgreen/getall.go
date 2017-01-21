package main

import "github.com/urfave/cli"

func doGetAll(c *cli.Context) error {
	client := getClient(c.String("addr"))
	entries, err := client.GetAll()
	if err != nil {
		return err
	}
	return printEntries(entries)
}

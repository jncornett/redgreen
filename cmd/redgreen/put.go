package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jncornett/redgreen"
	"github.com/jncornett/restful"
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
	var lines []string
	if c.Bool("stdin") {
		lines, err = readAllLines(os.Stdin)
		if err != nil {
			return err
		}
	}
	client := getClient(c.String("addr"))
	_, err = client.Put(redgreen.Entry{
		ID:   restful.ID(key),
		OK:   ok,
		Data: lines,
	})
	return err
}

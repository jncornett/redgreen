package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/jncornett/redgreen"
	"github.com/jncornett/restful"
	"github.com/urfave/cli"
)

//go:generate go-bindata data/...

const (
	defaultListenAddr  = ":8080"
	defaultConnectAddr = "localhost" + defaultListenAddr
	apiEndpoint        = "/api"
)

var (
	serverFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Value: defaultListenAddr,
			Usage: "address to listen to",
		},
	}
	clientFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Value: defaultConnectAddr,
			Usage: "address to connect to",
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "start a server",
			Action:  doServe,
			Flags:   serverFlags,
		},
		{
			Name:           "run",
			Aliases:        []string{"r", "x"},
			Usage:          "run a shell command",
			Action:         doRun,
			SkipArgReorder: true, // so we don't parse command flags
			Flags: append([]cli.Flag{cli.BoolFlag{
				Name:  "redir",
				Usage: "tee STDERR to STDOUT",
			}}, clientFlags...),
		},
		{
			Name:     "put",
			Aliases:  []string{"p"},
			Usage:    "insert an entry",
			Action:   doPut,
			Category: "op",
			Flags: append([]cli.Flag{cli.BoolFlag{
				Name:  "stdin",
				Usage: "read data from STDIN",
			}}, clientFlags...),
		},
		{
			Name:     "get",
			Aliases:  []string{"g"},
			Usage:    "lookup an entry",
			Action:   doGet,
			Category: "op",
			Flags:    clientFlags,
		},
		{
			Name:     "delete",
			Aliases:  []string{"d", "del"},
			Usage:    "lookup and delete an entry",
			Action:   doDelete,
			Category: "op",
			Flags:    clientFlags,
		},
		{
			Name:     "getall",
			Aliases:  []string{"ga"},
			Usage:    "lookup all entries",
			Action:   doGetAll,
			Category: "op",
			Flags:    clientFlags,
		},
	}
	app.Run(os.Args)
}

func getClient(addr string) restful.ClientStore {
	return restful.NewJSONClient(
		"http://"+addr+apiEndpoint,
		func() interface{} { return &redgreen.Entry{} },
		func() interface{} { return &[]redgreen.Entry{} },
	)
}

func printEntries(entries []redgreen.Entry) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func readAllLines(r io.Reader) (out []string, err error) {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		out = append(out, s.Text())
	}
	err = s.Err()
	return
}

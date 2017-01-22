package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/jncornett/redgreen"
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
			Aliases:  []string{"pu"},
			Usage:    "insert an entry",
			Action:   doPut,
			Category: "op",
			Flags:    clientFlags,
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
			Name:     "pop",
			Aliases:  []string{"po"},
			Usage:    "lookup and delete an entry",
			Action:   doPop,
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

func getClient(addr string) redgreen.HTTPJSONClient {
	return redgreen.HTTPJSONClient("http://" + addr + apiEndpoint)
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

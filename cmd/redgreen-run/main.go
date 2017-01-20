package main

import (
	"flag"
	"log"
	"os/exec"

	"github.com/jncornett/redgreen"
)

const defaultAddr = "localhost:8080"

func main() {
	var (
		addr = flag.String("addr", defaultAddr, "address of redgreen server")
	)
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("not enough arguments: must have key command [args]...")
	}
	key := args[0]
	cmd := exec.Command(args[1], args[2:]...)
	err := cmd.Run()
	rg := redgreen.RedGreenHTTPJSONClient(*addr)
	rg.Put(redgreen.Entry{Key: key, Value: err == nil})
	if err != nil {
		log.Fatal(err)
	}
}

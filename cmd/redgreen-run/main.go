package main

import (
	"flag"
	"log"
)

const defaultAddr = "localhost:8080"

func main() {
	var (
		addr = flag.String("addr", defaultAddr, "address of redgreen server")
	)
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("not enough arguments: must have key [command]...")
	}
	_ := args[0]
	_ := args[1:]
}

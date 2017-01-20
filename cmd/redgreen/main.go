package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/jncornett/redgreen"
)

//go:generate go-bindata data/...

const (
	defaultListenAddr = ":8080"
)

func logRequest(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		h.ServeHTTP(w, r)
	}
}

func logRequestFunc(h http.HandlerFunc) http.HandlerFunc {
	return logRequest(h)
}

func cli(addr string, args []string) {
	rg := redgreen.RedGreenHTTPJSONClient(addr)
	if len(args) < 1 {
		log.Fatal("must have one of {put,get} for first arg in CLI mode")
	}
	cmd := args[0]
	if cmd == "put" {
		if len(args) < 3 {
			log.Fatal("must specify put [key] [value], not enough args")
		}
		key := args[1]
		var value bool
		if args[2] == "true" {
			value = true
		} else if args[2] != "false" {
			log.Fatalln("value must be one of {true,false}, not", args[2])
		}
		e := redgreen.Entry{Key: key, Value: value}
		rg.Put(e)
	} else if cmd == "get" {
		if len(args) >= 2 {
			fmt.Println(rg.Get(args[1]))
		} else {
			fmt.Println(rg.GetAll())
		}
	} else {
		log.Fatal("must have one of {put,get} for first arg in CLI mode")
	}
}

func serve(addr string) {
	rg := redgreen.NewRedGreenMaster()
	defer rg.Close()
	http.Handle("/api", logRequest(http.StripPrefix("/api", redgreen.NewRedGreenHTTPJSONServer(rg))))
	http.Handle("/", logRequest(http.FileServer(&assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "data/static",
	})))
	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	var (
		listenAddr = flag.String("listen", defaultListenAddr, "address to listen on")
		doServe    = flag.Bool("serve", false, "toggle between client and server")
	)
	flag.Parse()
	if *doServe {
		serve(*listenAddr)
	} else {
		if *listenAddr == defaultListenAddr {
			*listenAddr = "localhost" + defaultListenAddr
		}
		cli(*listenAddr, flag.Args())
	}
}

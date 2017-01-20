package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
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
		e := Entry{Key: key, Value: value}
		var b bytes.Buffer
		json.NewEncoder(&b).Encode(e)
		resp, err := http.Post(addr+"/api", "application/json; charset=utf-8", &b)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(resp)
	} else if cmd == "get" {
		var (
			resp *http.Response
			err  error
		)
		if len(args) >= 2 {
			resp, err = http.Get(addr + "/api/" + args[1])
		} else {
			resp, err = http.Get(addr + "/api")
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp)
	} else {
		log.Fatal("must have one of {put,get} for first arg in CLI mode")
	}
}

func serve(addr string) {
	rg := NewRedGreen()
	defer rg.Close()
	http.HandleFunc("/api", logRequestFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(rg.GetAll())
		} else if r.Method == "POST" {
			var e Entry
			err := json.NewDecoder(r.Body).Decode(&e)
			if err != nil {
				http.Error(w, err.Error(), 400)
			}
			rg.Put(e)
		}
	}))
	http.Handle("/api/", http.StripPrefix("/api/", logRequestFunc(func(w http.ResponseWriter, r *http.Request) {
		var e Entry
		if r.Body == nil {
			http.Error(w, "no request body", 400)
			return
		}
		if r.Method != "GET" {
			http.Error(w, fmt.Sprintln("invalid method:", r.Method), 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&e)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		e = rg.Get(e.Key)
		json.NewEncoder(w).Encode(e)
	})))
	http.Handle("/", http.FileServer(&assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "data/static",
	}))
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
			*listenAddr = "http://" + "localhost" + defaultListenAddr
		}
		cli(*listenAddr, flag.Args())
	}
}

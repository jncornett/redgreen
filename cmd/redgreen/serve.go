package main

import (
	"fmt"
	"log"
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/jncornett/redgreen"
	"github.com/urfave/cli"
)

func doServe(c *cli.Context) error {
	store := redgreen.NewStore()
	defer store.Close()
	http.Handle(
		apiEndpoint,
		http.StripPrefix(apiEndpoint, logger(redgreen.NewHTTPJSONServer(store))),
	)
	http.Handle(
		apiEndpoint+"/",
		http.StripPrefix(apiEndpoint+"/", logger(redgreen.NewHTTPJSONServer(store))),
	)
	http.Handle("/", http.FileServer(&assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "data/static",
	}))
	return http.ListenAndServe(c.String("addr"), nil)
}

func logger(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		var resp string
		if len(r.URL.Path) != 0 {
			resp = fmt.Sprintf("- %v", r.URL)
		}
		log.Println(r.Method, resp)
	}
}

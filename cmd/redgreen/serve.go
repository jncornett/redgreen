package main

import (
	"fmt"
	"log"
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/jncornett/overlayfs"
	"github.com/jncornett/redgreen"
	"github.com/jncornett/restful"
	"github.com/urfave/cli"
)

func doServe(c *cli.Context) error {
	store := redgreen.NewStore()
	handler := http.StripPrefix(apiEndpoint, logger(restful.NewJSONHandler(store)))
	http.Handle(apiEndpoint, handler)
	http.Handle(apiEndpoint+"/", handler)
	http.Handle("/", http.FileServer(getFs(c)))
	return http.ListenAndServe(c.String("addr"), nil)
}

func getFs(c *cli.Context) http.FileSystem {
	assetFs := &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "data/static",
	}
	overlay := c.String("static")
	if overlay == "" {
		return assetFs
	}
	return overlayfs.NewOverlayFs(
		http.Dir(overlay),
		assetFs,
	)
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

package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/jncornett/redgreen"
	"github.com/jncornett/restful"
	"github.com/urfave/cli"
)

func doMock(c *cli.Context) error {
	store := redgreen.NewStore()
	handler := http.StripPrefix(apiEndpoint, logger(restful.NewJSONHandler(store)))
	http.Handle(apiEndpoint, handler)
	http.Handle(apiEndpoint+"/", handler)
	http.Handle("/", http.FileServer(&assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "data/static",
	}))
	go randomMockUpdates(store)
	return http.ListenAndServe(defaultListenAddr, nil)
}

func randomMockUpdates(store restful.Store) {
	for {
		waitS := time.Duration(rand.Intn(10)+5) * time.Second
		log.Println("CHAOSMONKEY -- sleep for", waitS)
		time.Sleep(waitS)
		log.Println("CHAOSMONKEY -- ACT!")
		v, err := store.GetAll()
		if err != nil {
			continue
		}
		all := v.([]redgreen.Entry)
		exist := choose(all)
		if exist == nil || decide(50) {
			// put
			entry := &redgreen.Entry{
				OK:   decide(50),
				Data: randomStrings(),
			}
			if exist == nil || decide(50) {
				// create
				entry.ID = restful.ID(randomString())
				if _, err := store.Put(entry); err != nil {
					log.Println(err)
				}
			} else {
				// update
				entry.ID = exist.ID
				if err := store.Update(entry.ID, entry); err != nil {
					log.Println(err)
				}
			}
		} else {
			// delete
			if err := store.Delete(exist.ID); err != nil {
				log.Println(err)
			}
		}
	}
}

func decide(p int) bool {
	return rand.Intn(100) > p
}

func choose(entries []redgreen.Entry) *redgreen.Entry {
	if len(entries) == 0 {
		return nil
	}
	return &entries[rand.Intn(len(entries))]
}

func randomString() string {
	runes := []rune(`abcdefghijklmnopqrstuvwxyz0123456789`)
	x := rand.Intn(20) + 5
	var out []rune
	for i := 0; i < x; i++ {
		out = append(out, runes[rand.Intn(len(runes))])
	}
	return string(out)
}

func randomStrings() []string {
	n := rand.Intn(5)
	var out []string
	for i := 0; i < n; i++ {
		out = append(out, randomString())
	}
	return out
}

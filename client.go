package redgreen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type RedGreenHTTPJSONClient string

func (rg *RedGreenHTTPJSONClient) GetAll() []Entry {
	resp, err := http.Get(rg.addr("api"))
	if err != nil {
		log.Fatal(err)
	}
	var entries []Entry
	err = json.NewDecoder(resp.Body).Decode(&entries)
	if err != nil {
		log.Fatal(err)
	}
	return entries
}

func (rg *RedGreenHTTPJSONClient) Get(k string) Entry {
	resp, err := http.Get(rg.addr("api", k))
	if err != nil {
		log.Fatal(err)
	}
	var entry Entry
	err = json.NewDecoder(resp.Body).Decode(&entry)
	if err != nil {
		log.Fatal(err)
	}
	return entry
}

func (rg *RedGreenHTTPJSONClient) Put(e Entry) {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(e)
	_, err := http.Post(rg.addr("api"), "application/json; charset=utf-8", &b)
	if err != nil {
		log.Fatal(err)
	}
}

func (rg RedGreenHTTPJSONClient) Close() {}

func (rg RedGreenHTTPJSONClient) addr(endpoints ...string) string {
	return fmt.Sprintf("http://%v/%v", rg, strings.Join(endpoints, "/"))
}

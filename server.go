package redgreen

import (
	"encoding/json"
	"net/http"
	"strings"
)

func NewRedGreenHTTPJSONServer(rg RedGreen) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.URL.Path, "/")
		if key == "" {
			if r.Method == "GET" {
				handleGetAll(rg, w, r)
			} else if r.Method == "POST" {
				handlePut(rg, w, r)
			}
		} else {
			handleGet(rg, key, w, r)
		}
	}
}

func handleGetAll(rg RedGreen, w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(rg.GetAll())
}

func handlePut(rg RedGreen, w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "no request body", 400)
		return
	}
	var entry Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	rg.Put(entry)
}

func handleGet(rg RedGreen, key string, w http.ResponseWriter, r *http.Request) {
	entry := rg.Get(key)
	json.NewEncoder(w).Encode(&entry)
}

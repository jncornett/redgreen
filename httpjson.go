package redgreen

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

type HTTPJSONClient string

func (c HTTPJSONClient) Get(key string) (entry Entry, ok bool, err error) {
	var resp *http.Response
	resp, err = http.Get(c.getEndpoint(key))
	if err != nil {
		return
	}
	ok = resp.StatusCode != http.StatusNotFound
	if ok {
		err = json.NewDecoder(resp.Body).Decode(&entry)
	}
	return
}

func (c HTTPJSONClient) GetAll() (out []Entry, err error) {
	var resp *http.Response
	resp, err = http.Get(string(c))
	if err != nil {
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&out)
	return
}

func (c HTTPJSONClient) Put(entry Entry) (err error) {
	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(&entry)
	if err != nil {
		return
	}
	_, err = http.Post(string(c), "application/json; charset=utf-8", &b)
	return
}

func (c HTTPJSONClient) Pop(key string) (entry Entry, ok bool, err error) {
	var req *http.Request
	req, err = http.NewRequest("DELETE", c.getEndpoint(key), nil)
	if err != nil {
		return
	}
	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	ok = resp.StatusCode != http.StatusNotFound
	if ok {
		err = json.NewDecoder(resp.Body).Decode(&entry)
	}
	return
}

func (c HTTPJSONClient) getEndpoint(parts ...string) string {
	return strings.Join([]string{string(c), strings.Join(parts, "/")}, "/")
}

var _ Client = HTTPJSONClient("")

type HTTPJSONServer struct {
	store Store
}

func NewHTTPJSONServer(s Store) *HTTPJSONServer {
	return &HTTPJSONServer{store: s}
}

func (s HTTPJSONServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if len(r.URL.Path) == 0 {
			s.handleGetAll(w, r)
		} else {
			s.handleGet(w, r)
		}
	case "POST":
		s.handlePut(w, r)
	case "DELETE":
		s.handlePop(w, r)
	}
}

func (s HTTPJSONServer) handleGet(w http.ResponseWriter, r *http.Request) {
	entry, ok := s.store.Get(r.URL.Path)
	if !ok {
		http.Error(w, "", http.StatusNotFound) // TODO add error string
		return
	}
	err := json.NewEncoder(w).Encode(&entry)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError) // TODO add error string
		return
	}
}

func (s HTTPJSONServer) handleGetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.store.GetAll())
}

func (s HTTPJSONServer) handlePut(w http.ResponseWriter, r *http.Request) {
	var entry Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // TODO add context to error string
		return
	}
	s.store.Put(entry)
}

func (s HTTPJSONServer) handlePop(w http.ResponseWriter, r *http.Request) {
	entry, ok := s.store.Pop(r.URL.Path)
	if !ok {
		http.Error(w, "", http.StatusNotFound) // TODO add error string
		return
	}
	err := json.NewEncoder(w).Encode(&entry)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError) // TODO add error string
		return
	}
}

var _ http.Handler = &HTTPJSONServer{}

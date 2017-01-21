package redgreen

import "time"

type MapStore map[string]Entry

func NewMapStore() MapStore {
	return make(MapStore)
}

func (s MapStore) Get(key string) (entry Entry, ok bool) {
	entry, ok = s[key]
	return
}

func (s MapStore) GetAll() (out []Entry) {
	for _, entry := range s {
		out = append(out, entry)
	}
	return
}

func (s MapStore) Put(entry Entry) {
	if entry.Updated.IsZero() {
		entry.Updated = time.Now()
	}
	s[entry.Key] = entry
}

func (s MapStore) Pop(key string) (entry Entry, ok bool) {
	entry, ok = s[key]
	delete(s, key)
	return
}

var _ Store = make(MapStore)

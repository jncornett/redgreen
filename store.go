package redgreen

import "time"

type Entry struct {
	Key     string    `json:"key"`
	OK      bool      `json:"ok"`
	Data    []string  `json:"data"`
	Updated time.Time `json:"updated"`
}

type Store interface {
	Get(string) (Entry, bool)
	GetAll() []Entry
	Put(Entry)
	Pop(string) (Entry, bool)
}

func NewStore() *ConcurrentStore {
	return NewConcurrentStore(NewRWStore(NewMapStore()))
}

package main

type Entry struct {
	Key   string `json:"key"`
	Value bool   `json:"value"`
}

type RedGreen interface {
	GetAll() []Entry
	Get(string) Entry
	Put(Entry)
	Close()
}

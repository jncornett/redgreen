package main

import "sync"

const RedGreenChanBound = 15

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

type redGreen struct {
	in   chan Entry
	data map[string]bool
	rw   sync.RWMutex
}

func NewRedGreen() RedGreen {
	rg := &redGreen{
		in:   make(chan Entry, RedGreenChanBound),
		data: make(map[string]bool),
	}
	rg.start()
	return rg
}

func (rg *redGreen) GetAll() (out []Entry) {
	rg.rw.RLock()
	defer rg.rw.RUnlock()
	for k, v := range rg.data {
		out = append(out, Entry{Key: k, Value: v})
	}
	return
}

func (rg *redGreen) Get(k string) Entry {
	rg.rw.RLock()
	defer rg.rw.RUnlock()
	return Entry{Key: k, Value: rg.data[k]}
}

func (rg *redGreen) Put(e Entry) {
	rg.in <- e
}

func (rg *redGreen) Close() {
	close(rg.in)
}

func (rg *redGreen) put(e Entry) {
	rg.rw.Lock()
	defer rg.rw.Unlock()
	rg.data[e.Key] = e.Value
}

func (rg *redGreen) start() {
	go func() {
		for {
			e, ok := <-rg.in
			if !ok {
				break
			}
			rg.put(e)
		}
	}()
}

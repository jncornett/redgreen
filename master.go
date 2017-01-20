package main

import "sync"

const rgMasterChannelBound = 15

type RedGreenMaster struct {
	in   chan Entry
	data map[string]bool
	rw   sync.RWMutex
}

func NewRedGreenMaster() RedGreen {
	rg := &RedGreenMaster{
		in:   make(chan Entry, rgMasterChannelBound),
		data: make(map[string]bool),
	}
	rg.start()
	return rg
}

func (rg *RedGreenMaster) GetAll() (out []Entry) {
	rg.rw.RLock()
	defer rg.rw.RUnlock()
	for k, v := range rg.data {
		out = append(out, Entry{Key: k, Value: v})
	}
	return
}

func (rg *RedGreenMaster) Get(k string) Entry {
	rg.rw.RLock()
	defer rg.rw.RUnlock()
	return Entry{Key: k, Value: rg.data[k]}
}

func (rg *RedGreenMaster) Put(e Entry) {
	rg.in <- e
}

func (rg *RedGreenMaster) Close() {
	close(rg.in)
}

func (rg *RedGreenMaster) put(e Entry) {
	rg.rw.Lock()
	defer rg.rw.Unlock()
	rg.data[e.Key] = e.Value
}

func (rg *RedGreenMaster) start() {
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

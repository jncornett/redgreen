package redgreen

import (
	"errors"
	"time"

	"github.com/jncornett/restful"
)

type MapStore map[restful.ID]Entry

func NewMapStore() MapStore { return make(MapStore) }

func (s MapStore) Get(key restful.ID) (interface{}, error) {
	entry, ok := s[key]
	if !ok {
		return nil, restful.ErrMissing{key}
	}
	return entry, nil
}

func (s MapStore) GetAll() (interface{}, error) {
	entries := []Entry{}
	for _, entry := range s {
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s MapStore) Put(v interface{}) (interface{}, error) {
	entry, ok := v.(*Entry)
	if !ok {
		return nil, errors.New("wrong type") // FIXME could be panic
	}
	if entry.Updated.IsZero() {
		entry.Updated = time.Now()
	}
	s[entry.ID] = *entry
	return &entry, nil
}

func (s MapStore) Update(id restful.ID, v interface{}) error {
	if _, ok := s[id]; !ok {
		return restful.ErrMissing{id}
	}
	entry, ok := v.(*Entry)
	if !ok {
		return errors.New("wrong type") // FIXME could be panic
	}
	if entry.Updated.IsZero() {
		entry.Updated = time.Now()
	}
	entry.ID = id // make sure ID is in sync
	s[entry.ID] = *entry
	return nil
}

func (s MapStore) Delete(key restful.ID) error {
	if _, ok := s[key]; !ok {
		return restful.ErrMissing{key}
	}
	delete(s, key)
	return nil
}

func (s MapStore) New() interface{} {
	return &Entry{}
}

var _ restful.Store = NewMapStore()

func NewStore() restful.Store {
	return restful.NewRWStore(NewMapStore())
}

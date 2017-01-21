package redgreen

import "sync"

type RWStore struct {
	store        Store
	sync.RWMutex // FIXME don't want mux methods accessible from outside
}

func NewRWStore(s Store) *RWStore {
	return &RWStore{store: s}
}

func (s *RWStore) Get(key string) (Entry, bool) {
	s.RLock()
	defer s.RUnlock()
	return s.store.Get(key)
}

func (s *RWStore) GetAll() []Entry {
	s.RLock()
	defer s.RUnlock()
	return s.store.GetAll()
}

func (s *RWStore) Put(entry Entry) {
	s.Lock()
	defer s.Unlock()
	s.store.Put(entry)
}

func (s *RWStore) Pop(key string) (Entry, bool) {
	s.Lock()
	defer s.Unlock()
	return s.store.Pop(key)
}

var _ Store = &RWStore{}

const ConcurrentWriteBd = 5

type storeOp struct {
	entry Entry
	pop   bool
}

type ConcurrentStore struct {
	store Store
	q     chan storeOp
}

func NewConcurrentStore(s Store) *ConcurrentStore {
	store := &ConcurrentStore{
		store: s,
		q:     make(chan storeOp, ConcurrentWriteBd),
	}
	go func() {
		for {
			op, ok := <-store.q
			if !ok {
				return // channel closed, we're done
			}
			if op.pop {
				store.store.Pop(op.entry.Key)
			} else {
				store.store.Put(op.entry)
			}
		}
	}()
	return store
}

func (s *ConcurrentStore) Get(key string) (Entry, bool) {
	return s.store.Get(key)
}

func (s *ConcurrentStore) GetAll() []Entry {
	return s.store.GetAll()
}

func (s *ConcurrentStore) Put(entry Entry) {
	s.q <- storeOp{entry, false}
}

// NOTE Pop is nonatomic for ConcurrentStore
func (s *ConcurrentStore) Pop(key string) (entry Entry, ok bool) {
	entry, ok = s.store.Get(key)
	if ok {
		s.q <- storeOp{entry, true}
	}
	return
}

func (s *ConcurrentStore) Close() {
	close(s.q)
}

var _ Store = &ConcurrentStore{}

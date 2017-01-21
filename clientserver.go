package redgreen

type Client interface {
	Get(string) (Entry, bool, error)
	GetAll() ([]Entry, error)
	Put(Entry) error
	Pop(string) (Entry, bool, error)
}

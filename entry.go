package redgreen

import (
	"time"

	"github.com/jncornett/restful"
)

type Entry struct {
	ID      restful.ID `json:"key"`
	OK      bool       `json:"ok"`
	Data    []string   `json:"data"`
	Updated time.Time  `json:"updated"`
}

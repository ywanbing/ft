package msg

import (
	"sync"
)

var (
	msgPool = sync.Pool{
		New: func() any {
			return &Message{}
		},
	}
)

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

	BytesPool = sync.Pool{
		New: func() any {
			return make([]byte, 40*1024)
		},
	}
)

package msg

import (
	"strings"
	"sync"
)

var (
	msgPool = sync.Pool{
		New: func() any {
			return &Message{}
		},
	}

	builderPool = sync.Pool{
		New: func() any {
			return &strings.Builder{}
		},
	}
)

package server

import (
	"fmt"
)

var (
	MAGIC_BYTES = []byte("f00t")
	EmErr       = fmt.Errorf("dont have msg")
)

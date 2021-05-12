package internal

import (
	"encoding/json"
	"github.com/google/uuid"
)

var MAGIC_BYTES = []byte("f00t")

func GenFileName() string {
	u := uuid.New()
	return u.String()

}

type Message struct {
	MsgType  MsgType `json:"t"`
	FileName string  `json:"f"`
	Bytes    []byte  `json:"b"`
	Size     uint64  `json:"s"`
}

func (m Message) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

// Decode will convert from bytes
func Decode(b []byte) (m Message, err error) {
	err = json.Unmarshal(b, &m)
	return
}

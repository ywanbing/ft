package msg

import (
	"encoding/json"
	"strconv"
	"strings"
)

type MsgType byte

const (
	MsgInvalid MsgType = iota
	MsgHead
	MsgFile
	MsgEnd
	MsgClose
)

type Message struct {
	MsgType  MsgType `json:"t"`
	FileName string  `json:"f"`
	Bytes    []byte  `json:"b"`
	Size     uint64  `json:"s"`
}

func (m *Message) GC() {
	m.reset()
	msgPool.Put(m)
}

func (m *Message) reset() {
	m.MsgType = MsgInvalid
	m.FileName = ""
	m.Bytes = nil
	m.Size = 0
}

func (m *Message) String() string {
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	builder.WriteString("{")

	builder.WriteString(`"t":`)
	builder.WriteString(strconv.Itoa(int(m.MsgType)) + ",")

	builder.WriteString(`"f":`)
	builder.WriteString(`"` + m.FileName + `",`)

	builder.WriteString(`"b":`)
	builder.WriteString(`"` + string(m.Bytes) + `",`)

	builder.WriteString(`"s":`)
	builder.WriteString(strconv.Itoa(int(m.Size)))

	builder.WriteString("}")
	return builder.String()
}

// Decode will convert from bytes
func Decode(b []byte) (m *Message, err error) {
	m = msgPool.Get().(*Message)
	err = json.Unmarshal(b, &m)
	return
}

package msg

import (
	"encoding/json"
)

type MsgType byte

const (
	MsgInvalid MsgType = iota
	MsgHead
	MsgFile
	MsgEnd
	MsgNotify
	MsgClose
)

type Status byte

const (
	Status_Ok Status = iota
	Status_Err
)

type Message struct {
	MsgType  MsgType `json:"t"`
	FileName string  `json:"f"`
	Bytes    []byte  `json:"b"`
	Size     uint64  `json:"s"`
}

type Notify struct {
	Status byte
}

func (m *Message) GC() {
	if m.MsgType == MsgFile {
		BytesPool.Put(m.Bytes[:cap(m.Bytes)])
	}
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
	bytes, _ := json.Marshal(m)
	return string(bytes)
}

// Decode will convert from bytes
func Decode(b []byte) (m *Message, err error) {
	m = msgPool.Get().(*Message)
	err = json.Unmarshal(b, &m)
	return
}

func NewNotifyMsg(fileName string, status Status) *Message {
	m := msgPool.Get().(*Message)
	m.MsgType = MsgNotify
	m.Bytes = []byte{byte(status)}
	m.FileName = fileName
	return m
}

func NewHeadMsg(fileName string) *Message {
	m := msgPool.Get().(*Message)
	m.MsgType = MsgHead
	m.FileName = fileName
	return m
}

func NewFileMsg(fileName string, buf []byte) *Message {
	m := msgPool.Get().(*Message)
	m.MsgType = MsgFile
	m.FileName = fileName
	m.Bytes = buf
	// m.Bytes = make([]byte, len(buf))
	// copy(m.Bytes, buf)
	return m
}

func NewEndMsg(fileName string, size uint64) *Message {
	m := msgPool.Get().(*Message)
	m.MsgType = MsgEnd
	m.FileName = fileName
	m.Size = size
	return m
}

func NewCloseMsg(fileName string, status Status) *Message {
	m := msgPool.Get().(*Message)
	m.MsgType = MsgClose
	m.Bytes = []byte{byte(status)}
	m.FileName = fileName
	return m
}

package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/ywanbing/ft/pkg/msg"
)

type TcpCon struct {
	conn *net.TCPConn

	recv chan *msg.Message
	send chan *msg.Message

	stop bool
}

func NewTcp(conn *net.TCPConn) *TcpCon {
	return &TcpCon{
		conn: conn,
		recv: make(chan *msg.Message, 1024),
		send: make(chan *msg.Message, 1024),
	}
}

func (t *TcpCon) HandlerLoop() {
	go t.readMsg()
	go t.sendMsg()
}

func (t *TcpCon) sendMsg() {
	defer t.conn.Close()

	var err error
	defer func() {
		if err != nil {
			fmt.Printf("found mistake: %s \n", err)
		}
	}()

	buf := make([]byte, 64*1024)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for !t.stop {
		select {
		case m := <-t.send:
			data := m.String()
			m.GC()

			dataLen := len(data)

			copy(buf[:4], MAGIC_BYTES)
			binary.BigEndian.PutUint32(buf[4:8], uint32(dataLen))
			copy(buf[8:], []byte(data))

			_, err = t.conn.Write(buf[:8+dataLen])
			if err != nil {
				return
			}
		case <-ticker.C:
			fmt.Println("wait send msg ... ")
		}
	}
}

func (t *TcpCon) readMsg() {
	defer t.conn.Close()

	var err error
	defer func() {
		if err != nil {
			fmt.Printf("found mistake: %s \n", err)
		}
		t.stop = true
	}()

	header := make([]byte, 4)
	buf := make([]byte, 64*1024)

	for {
		// read until we get 4 bytes for the magic
		_, err = io.ReadFull(t.conn, header)
		if err != nil {
			if err != io.EOF {
				err = fmt.Errorf("initial read error: %v \n", err)
				return
			}
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if !bytes.Equal(header, MAGIC_BYTES) {
			err = fmt.Errorf("initial bytes are not magic: %s", header)
			return
		}

		// read until we get 4 bytes for the header
		_, err = io.ReadFull(t.conn, header)
		if err != nil {
			err = fmt.Errorf("initial read error: %v \n", err)
			return
		}

		// 数据大小
		msgSize := binary.BigEndian.Uint32(header)

		// 解析为结构体消息
		var n int
		var m *msg.Message

		n, err = io.ReadFull(t.conn, buf[:msgSize])
		if err != nil {
			err = fmt.Errorf("initial read error: %v \n", err)
			return
		}

		m, err = msg.Decode(buf[:n])
		if err != nil {
			err = fmt.Errorf("read message error: %v \n", err)
			return
		}

		t.recv <- m
	}
}

func (t *TcpCon) GetMsg() (*msg.Message, bool) {
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	select {
	case m := <-t.recv:
		return m, true
	case <-timer.C:
		return nil, false
	}
}

func (t *TcpCon) SendMsg(m *msg.Message) {
	t.send <- m
}

func (t *TcpCon) Close() error {
	t.stop = true
	return nil
}

var _ = NetConn(&TcpCon{})

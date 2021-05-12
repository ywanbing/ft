package internal

import (
	"fmt"
	"net"
	"time"
)

var EmErr = fmt.Errorf("dont have msg")

type Client struct {
	c       *net.TCPConn
	receive chan Message
	send    chan Message
}

func NewClient(c *net.TCPConn) *Client {
	return &Client{
		c:       c,
		receive: make(chan Message, 512),
		send:    make(chan Message, 512),
	}
}

type Server struct {
	c       *net.TCPConn
	receive chan Message
	send    chan Message
}

func (s *Server) SendMsg(m Message) {
	s.send <- m
}

func (s *Server) Receive() (Message, error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	select {
	case m := <-s.receive:
		return m, nil
	case <-ticker.C:
		return Message{}, EmErr
	}
}

func (c *Client) SendMsg(m Message) {
	c.send <- m
}

func (c *Client) Receive() (Message, error) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	select {
	case m := <-c.receive:
		return m, nil
	case <-ticker.C:
		return Message{}, EmErr
	}
}

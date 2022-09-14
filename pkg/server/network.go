package server

import (
	"github.com/ywanbing/ft/pkg/msg"
)

type NetConn interface {
	// HandlerLoop 不能阻塞
	HandlerLoop()
	GetMsg() (*msg.Message, bool)
	SendMsg(m *msg.Message)
	Close() error
	IsClose() bool
}

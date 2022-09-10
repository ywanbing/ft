package file

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/ywanbing/ft/pkg/msg"
	"github.com/ywanbing/ft/pkg/server"
)

type Server interface {
	Start() error
}

type Client interface {
	SendFile() error
}

type ConMgr struct {
	conn     server.NetConn
	dir      string
	fileName string

	// 在发送重要消息的时候，需要同步等待消息的状态，返回是否正确
	waitNotify chan bool
	stop       bool
}

func NewServer(conn server.NetConn, dir string) Server {
	return &ConMgr{
		conn: conn,
		dir:  dir,
	}
}

func (c *ConMgr) Start() error {
	c.conn.HandlerLoop()
	// 处理接收的消息
	return c.handler()
}

func (c *ConMgr) handler() error {
	var fs *os.File
	var err error

	defer func() {
		if fs != nil {
			_ = fs.Close()
		}
		c.stop = true
	}()

	for !c.stop {
		m, ok := c.conn.GetMsg()
		if !ok {
			continue
		}

		switch m.MsgType {
		case msg.MsgHead:
			// 创建文件
			if m.FileName != "" {
				c.fileName = m.FileName
			} else {
				c.fileName = GenFileName()
			}

			fs, err = os.OpenFile(path.Clean(c.dir+"/"+c.fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				fmt.Println("os.Create err =", err)
				c.conn.SendMsg(msg.NewNotifyMsg(c.fileName, msg.Status_Err))
				return err
			}
			fmt.Println("send head is ok")

			c.conn.SendMsg(msg.NewNotifyMsg(c.fileName, msg.Status_Ok))
		case msg.MsgFile:
			if fs == nil {
				fmt.Println(c.fileName, "file is not open !")
				c.conn.SendMsg(msg.NewNotifyMsg(c.fileName, msg.Status_Err))
				return nil
			}
			// 写入文件
			_, err = fs.Write(m.Bytes)
			if err != nil {
				fmt.Println("file.Write err =", err)
				c.conn.SendMsg(msg.NewNotifyMsg(c.fileName, msg.Status_Err))
				return err
			}
		case msg.MsgEnd:
			// 操作完成
			info, _ := fs.Stat()
			if info.Size() != int64(m.Size) {
				err = fmt.Errorf("file.size %v rece size %v \n", info.Size(), m.Size)
				c.conn.SendMsg(msg.NewCloseMsg(c.fileName, msg.Status_Err))
				return err
			}

			fmt.Printf("save file %v is success \n", info.Name())
			c.conn.SendMsg(msg.NewNotifyMsg(c.fileName, msg.Status_Ok))
		case msg.MsgNotify:
			c.waitNotify <- m.Bytes[0] == byte(msg.Status_Ok)
		case msg.MsgClose:
			if m.Bytes[0] != byte(msg.Status_Ok) {
				return fmt.Errorf("server an error occurred")
			}
			return nil
		}
	}

	return err
}

func NewClient(conn server.NetConn, filePath string) Client {
	return &ConMgr{
		conn:       conn,
		fileName:   filePath,
		waitNotify: make(chan bool, 1),
	}
}

func (c *ConMgr) SendFile() error {
	var err error
	c.conn.HandlerLoop()
	// 处理接收的消息
	go func() {
		_ = c.handler()
	}()
	err = c.sendFile()
	return err
}

func (c *ConMgr) sendFile() error {
	defer func() {
		c.stop = true
		_ = c.conn.Close()
	}()

	file, err := os.Open(c.fileName)
	if err != nil {
		fmt.Printf("open file err %v \n", err)
		return err
	}
	fileInfo, _ := file.Stat()

	fmt.Println("client ready to write ...")
	m := msg.NewHeadMsg(fileInfo.Name())
	// 发送文件信息
	c.conn.SendMsg(m)

	// 等待服务器返回通知消息
	timer := time.NewTimer(5 * time.Second)
	select {
	case ok := <-c.waitNotify:
		if !ok {
			return fmt.Errorf("send err")
		}
	case <-timer.C:
		return fmt.Errorf("wait server msg timeout")
	}

	// 发送文件数据
	readBuf := make([]byte, 60*1024)
	for {
		n, err := file.Read(readBuf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		c.conn.SendMsg(msg.NewFileMsg(c.fileName, readBuf[:n]))
	}

	c.conn.SendMsg(msg.NewEndMsg(c.fileName, uint64(fileInfo.Size())))

	// 等待服务器返回通知消息
	timer = time.NewTimer(5 * time.Second)
	select {
	case ok := <-c.waitNotify:
		if !ok {
			return fmt.Errorf("send err")
		}
	case <-timer.C:
		return fmt.Errorf("wait server msg timeout")
	}

	fmt.Println("client send file success...")
	c.conn.SendMsg(msg.NewCloseMsg(c.fileName, msg.Status_Ok))
	return err
}

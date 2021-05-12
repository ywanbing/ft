package internal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"time"
)

type MsgType byte

const (
	MsgHead MsgType = 1
	MsgFile MsgType = 2
	MsgEnd  MsgType = 3
)

func StartServer(addr string, dir string) {
	//通过ResolveTCPAddr实例一个具体的tcp断点
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	//打开一个tcp断点监听
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()
	fmt.Println("Server ready to read ...")
	//循环接收客户端的连接，创建一个协程具体去处理连接
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		client := NewClient(tcpConn)
		fmt.Println("A client connected :" + tcpConn.RemoteAddr().String())
		go tcpPipe(client, dir)
	}
}

//具体处理连接过程方法
func tcpPipe(c *Client, dir string) {
	//tcp连接的地址
	ipStr := c.c.RemoteAddr().String()

	defer func() {
		fmt.Println(" Disconnected : " + ipStr)
		_ = c.c.Close()
	}()
	if !PathExists(dir) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
	Save := true
	go func() {
		defer func() {
			Save = false
		}()
		for {
			var err error
			// long read deadline in case waiting for file
			if err := c.c.SetReadDeadline(time.Now().Add(3 * time.Hour)); err != nil {
				fmt.Printf("error setting read deadline: %v \n", err)
			}
			// must clear the timeout setting
			defer c.c.SetDeadline(time.Time{})
			// read until we get 4 bytes for the magic
			header := make([]byte, 4)
			_, err = io.ReadFull(c.c, header)
			if err != nil {
				fmt.Printf("initial read error: %v \n", err)
				return
			}
			if !bytes.Equal(header, MAGIC_BYTES) {
				err = fmt.Errorf("initial bytes are not magic: %x", header)
				return
			}
			// read until we get 4 bytes for the header
			header = make([]byte, 4)
			_, err = io.ReadFull(c.c, header)
			if err != nil {
				fmt.Printf("initial read error: %v \n", err)
				return
			}
			numBytesUint32 := binary.BigEndian.Uint32(header)
			// shorten the reading deadline in case getting weird data
			if err := c.c.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
				fmt.Printf("error setting read deadline: %v \n", err)
			}
			buf := make([]byte, numBytesUint32)
			_, err = io.ReadFull(c.c, buf)
			if err != nil {
				fmt.Printf("consecutive read error: %v \n", err)
				return
			}
			m, err := Decode(buf)
			if err != nil {
				fmt.Printf("read message error: %v \n", err)
				return
			}
			c.receive <- m
		}
	}()

	fileName := GenFileName()
	var fs *os.File
	defer func() {
		if fs != nil {
			_ = fs.Close()
		}
	}()
	for {
		m, err := c.Receive()
		if err != nil {
			if !Save {
				break
			}
			fmt.Printf("receive err is %v \n", err)
			continue
		}

		switch m.MsgType {
		case MsgHead:
			// 创建文件
			if m.FileName != "" {
				fileName = m.FileName
			}
			fs, err = os.OpenFile(path.Clean(dir+"/"+fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				fmt.Println("os.Create err =", err)
				return
			}
			fmt.Println("send head is ok")
			_, _ = c.c.Write([]byte("ok"))
		case MsgFile:
			// 写入文件
			_, err := fs.Write(m.Bytes)
			if err != nil {
				fmt.Println("file.Write err =", err)
				return
			}
			_, _ = c.c.Write([]byte("ok"))
		case MsgEnd:
			// 操作完成
			info, _ := fs.Stat()
			if info.Size() != int64(m.Size) {
				fmt.Printf("file.size %v rece size %v \n", info.Size(), m.Size)
				return
			}
			fmt.Printf("save file %v is success \n", info.Name())
			_, _ = c.c.Write([]byte("end"))
		default:
			return
		}
	}
}

//PathExists 判断文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func StartClient(addr string, fileName string) (err error) {
	//通过ResolveTCPAddr实例一个具体的tcp断点
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	//打开一个tcp断点监听
	tcpListener, _ := net.DialTCP("tcp", nil, tcpAddr)
	defer tcpListener.Close()

	c := NewClient(tcpListener)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("open file err %v \n", err)
		return
	}
	fileInfo, _ := file.Stat()

	fmt.Println("client ready to write ...")
	m := Message{
		MsgType:  MsgHead,
		FileName: fileInfo.Name(),
		Bytes:    nil,
		Size:     0,
	}
	// 发送文件信息
	buf := []byte(m.String())

	data := make([]byte, len(MAGIC_BYTES)+4+len(buf))
	copy(data[:4], MAGIC_BYTES)
	binary.BigEndian.PutUint32(data[4:8], uint32(len(buf)))
	copy(data[8:], buf)

	_, err = c.c.Write(data)
	if err != nil {
		fmt.Println("conn.Write info.Name err =", err)
		return fmt.Errorf("send err")
	}
	buf = make([]byte, 1024)
	n, err := c.c.Read(buf)
	fmt.Printf("read msg is :[%v] \n", string(buf[:n]))
	if err != nil {
		return fmt.Errorf("send err")
	}
	if string(buf[:n]) != "ok" {
		return fmt.Errorf("send err")
	}
	// 发送文件数据

	readBuf := make([]byte, 60*1024)
	// 重新初始化
	data = make([]byte, 1*1024*1024)
	for {
		n, err := file.Read(readBuf)
		if err != nil && n == 0 {
			break
		}
		m.MsgType = MsgFile
		m.Bytes = readBuf[:n]
		buf = []byte(m.String())
		copy(data[:4], MAGIC_BYTES)
		binary.BigEndian.PutUint32(data[4:8], uint32(len(buf)))
		copy(data[8:], buf)
		_, err = c.c.Write(data[:8+len(buf)])
		if err != nil {
			break
		}
		readBuf := make([]byte, 1024)
		n, err = c.c.Read(readBuf)
		if err != nil {
			break
		}
		if string(readBuf[:n]) != "ok" {
			return fmt.Errorf("send err")
		}
	}

	m.MsgType = MsgEnd
	m.Bytes = []byte{}
	m.Size = uint64(fileInfo.Size())
	buf = []byte(m.String())

	binary.BigEndian.PutUint32(data[4:8], uint32(len(buf)))
	copy(data[8:], buf)

	_, err = c.c.Write(data[:8+len(buf)])
	if err != nil {
		return fmt.Errorf("send err")
	}
	readBuf = make([]byte, 1024)
	n, err = c.c.Read(readBuf)
	if err != nil {
		return err
	}
	fmt.Printf("read msg is [%v] \n", string(readBuf[:n]))
	if string(readBuf[:n]) != "end" {
		return fmt.Errorf("send err")
	}
	return nil
}

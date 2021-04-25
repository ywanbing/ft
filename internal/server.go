package internal

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

type MsgType byte

const (
	MsgHead MsgType = 1
	MsgFile MsgType = 2
	MsgEnd  MsgType = 3
)

// MsgHeader 1 + 7 +8(文件大小)
type MsgHeader [9]byte

func NewMsgHead(mt MsgType, size int64) MsgHeader {
	msg := [9]byte{}

	msg[0] = byte(mt)

	binary.BigEndian.PutUint64(msg[1:], uint64(size))

	return msg
}

func NewMsgHeader(mt MsgType, fileName string) []byte {

	f := []byte(fileName)
	head := NewMsgHead(mt, int64(len(f)))
	buf := make([]byte, len(head)+len(f))

	copy(buf, head[:])

	copy(buf[len(head):], f)

	return buf
}

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
		fmt.Println("A client connected :" + tcpConn.RemoteAddr().String())
		go tcpPipe(tcpConn, dir)
	}
}

//具体处理连接过程方法
func tcpPipe(conn *net.TCPConn, dir string) {
	//tcp连接的地址
	ipStr := conn.RemoteAddr().String()

	defer func() {
		fmt.Println(" Disconnected : " + ipStr)
		_ = conn.Close()
	}()
	if !PathExists(dir) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}

	buf := make([]byte, 9)

	n, err := conn.Read(buf)
	if err == io.EOF && n == 0 {
		return
	}

	if n != 9 {
		return
	}

	switch MsgType(buf[0]) {
	case MsgHead:
		// 创建文件
		size := binary.BigEndian.Uint64(buf[1:])
		data := make([]byte, size)
		n, _ = conn.Read(data)
		if uint64(n) != size {
			return
		}
		fileName := string(data)
		f := dir + "/" + fileName
		file, _ := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		defer file.Close()
		for {

		}



	case MsgFile:
		fallthrough
	case MsgEnd:
		_, _ = conn.Write([]byte("send data err "))
		return
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

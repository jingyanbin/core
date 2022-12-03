package xnet

import "net"

const FLAG_ENCRYPT byte = 1 << 0  //加密
const FLAG_COMPRESS byte = 1 << 1 //压缩

const netHeaderSize = 7
const netNeedCompressLength = 1500

type Conn interface {
	Close() error
	RemoteAddr() string
	Recv() ([]byte, bool)
	Send(msg []byte) bool
	HandshakeSend(flag byte) bool
	HandshakeRecvSeed() bool
}

func Read(conn net.Conn, buf []byte, size int) (n int, err error) {
	var nn int
	for n < size && err == nil {
		nn, err = conn.Read(buf[n:])
		n += nn
	}
	return
}

func Write(conn net.Conn, buf []byte, size int) (n int, err error) {
	var nn int
	for n < size && err == nil {
		nn, err = conn.Write(buf[n:])
		n += nn
	}
	return
}

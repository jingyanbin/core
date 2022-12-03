package xnet

import (
	"encoding/binary"
	"github.com/jingyanbin/core/internal"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
)

func NewTCPListener(address string) (lis *net.TCPListener, err error) {
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	lis, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	return lis, nil
}

func ConnectTCP(addr string, timeout time.Duration) (*net.TCPConn, error) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		internal.Log.Error("tcp client connect server failed: %s, %s", addr, err.Error())
		return nil, err
	}
	return conn.(*net.TCPConn), nil
}

type TCPConn struct {
	conn   *net.TCPConn
	crypt  XORCrypt
	head   []byte
	flag   byte
	closed int32
}

func (m *TCPConn) read(buf []byte, size int) (n int, err error) {
	var nn int
	for n < size && err == nil {
		nn, err = m.conn.Read(buf[n:])
		n += nn
	}
	return
}

func (m *TCPConn) write(buf []byte, size int) (n int, err error) {
	var nn int
	for n < size && err == nil {
		nn, err = m.conn.Write(buf[n:])
		n += nn
	}
	return
}

func (m *TCPConn) Close() error {
	if atomic.CompareAndSwapInt32(&m.closed, 0, 1) {
		return m.conn.Close()
	}
	return nil
}

func (m *TCPConn) RemoteAddr() string {
	return m.conn.RemoteAddr().String()
}

func (m *TCPConn) Recv() ([]byte, bool) {
	err := m.conn.SetReadDeadline(time.Now().Add(option.RecvTimeout))
	if err != nil {
		internal.Log.Error("tcp recv set read deadline error: %s, %s", m.RemoteAddr(), err.Error())
		return nil, false
	}
	var n int
	n, err = m.read(m.head, netHeaderSize)
	if err != nil {
		internal.Log.Error("tcp recv head error: %s, %d/%d, %s", m.RemoteAddr(), n, netHeaderSize, err.Error())
		return nil, false
	}
	flag, bodyBcc, headBcc := m.head[4], m.head[5], m.head[6]
	isEncrypt := flag&FLAG_ENCRYPT == FLAG_ENCRYPT
	if isEncrypt {
		if bcc := CountBCC(m.head, 0, 6); bcc != headBcc { //校验head
			internal.Log.Error("tcp recv conn head bcc error: %s, bcc=%v, head bcc=%v", m.RemoteAddr(), bcc, headBcc)
			return nil, false
		}
	}
	length := int(binary.BigEndian.Uint32(m.head[:4]))
	content := make([]byte, length)
	n, err = m.read(content, length)
	if err != nil {
		internal.Log.Error("tcp recv content error: %s, %d/%d, %s", m.RemoteAddr(), n, length, err.Error())
		return nil, false
	}
	if isEncrypt {
		if bcc := CountBCC(content, 0, length); bcc != bodyBcc {
			internal.Log.Error("tcp recv conn body bcc error: %s, bcc=%v, body bcc=%v", m.RemoteAddr(), bcc, bodyBcc)
			return nil, false
		}
		m.crypt.Decrypt(content, 0, length)
	}
	if flag&FLAG_COMPRESS == FLAG_COMPRESS {
		content, err = internal.Compress.UnGZip(content)
		if err != nil {
			internal.Log.Error("tcp recv conn un compress error: %s, %s", m.RemoteAddr(), err.Error())
			return nil, false
		}
	}
	return content, true
}

func (m *TCPConn) Send(msg []byte) bool {
	var length int
	length = len(msg)
	var flag byte
	if m.flag&FLAG_COMPRESS == FLAG_COMPRESS {
		if length > netNeedCompressLength {
			flag |= FLAG_COMPRESS
			msg = internal.Compress.GZip(msg)
			length = len(msg)
		}
	}
	max := netHeaderSize + length
	content := make([]byte, max)
	binary.BigEndian.PutUint32(content, uint32(length)) //length
	copy(content[7:], msg)
	if m.flag&FLAG_ENCRYPT == FLAG_ENCRYPT {
		content[4] = flag | FLAG_ENCRYPT //flag
		m.crypt.Encrypt(content, 7, length)
		content[5] = CountBCC(content, 7, length) //body bcc
		content[6] = CountBCC(content, 0, 6)      // head bcc
	} else {
		content[4] = flag
	}

	n, err := m.write(content, max)
	if err != nil {
		internal.Log.Error("tcp send error: %s, %d/%d, %s", m.RemoteAddr(), n, max, err.Error())
		return false
	}
	return true
}

func (m *TCPConn) HandshakeSend(flag byte) bool {
	m.flag = flag
	buf := make([]byte, 9)
	if flag&FLAG_ENCRYPT == FLAG_ENCRYPT {
		m.crypt.iSeed = rand.Uint32()
		m.crypt.oSeed = rand.Uint32()
		binary.BigEndian.PutUint32(buf[1:], m.crypt.iSeed)
		binary.BigEndian.PutUint32(buf[5:], m.crypt.oSeed)
		XOREncrypt(option.HandshakeSeed, buf, 0, 9)
	}
	n, err := m.write(buf, 9)
	if err != nil {
		internal.Log.Error("tcp handshake send seed write error: %s, %d, %s", m.RemoteAddr(), n, err.Error())
		return false
	}
	internal.Log.Error("tcp handshake send seed crypt: %v", m.crypt)
	return true
}

func (m *TCPConn) HandshakeRecvSeed() bool {
	buf := make([]byte, 9)
	err := m.conn.SetReadDeadline(time.Now().Add(option.HandshakeRecvTimeout)) //超时未收到握手信息握手失败
	if err != nil {
		internal.Log.Error("tcp handshake recv seed read deadline error: %s, %s", m.RemoteAddr(), err.Error())
		return false
	}
	n, err := m.read(buf, 9)
	if err != nil {
		internal.Log.Error("tcp handshake recv seed read error: %s, %d, %s", m.RemoteAddr(), n, err.Error())
		return false
	}
	XORDecrypt(option.HandshakeSeed, buf, 0, 9)
	m.flag = buf[0]
	m.crypt.iSeed = binary.BigEndian.Uint32(buf[5:])
	m.crypt.oSeed = binary.BigEndian.Uint32(buf[1:])
	internal.Log.Error("tcp handshake recv seed crypt: %v", m.crypt)
	return true
}

func NewTCPConn(conn *net.TCPConn) *TCPConn {
	return &TCPConn{conn: conn, head: make([]byte, netHeaderSize)}
}

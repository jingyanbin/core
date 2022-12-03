package xnet

import (
	"time"
)

type Option struct {
	RecvTimeout          time.Duration //接收超时
	HandshakeRecvTimeout time.Duration //握手接收超时
	HandshakeSeed        uint32        //握手加密种子
	CompressLen          uint32        //压缩长度
	XORCryptA            uint32        //参数A
	XORCryptB            uint32        //参数B
	XORBcc               byte          //初始BCC
}

var option = Option{
	RecvTimeout:          time.Second * 30,
	HandshakeRecvTimeout: time.Second * 5,
	HandshakeSeed:        3216,
	CompressLen:          1500,
	XORCryptA:            825,
	XORCryptB:            658,
	XORBcc:               17,
}

func Init(opt Option) {
	option = opt
}

package filequeue

import (
	"github.com/jingyanbin/core/basal"
	internal "github.com/jingyanbin/core/internal"
)

type Options struct {
	ConfDataDir      string //配置文件目录
	MsgFileDir       string //消息文件目录
	FileNamePrefix   string //文件名前缀
	Sep              []byte //消息数据分割符
	MsgFileMaxByte   int64  //pusher配置 消息文件最大占用空间 单位: 字节
	PushChanSize     int    //pusher配置 入队管道大小
	DeletePoppedFile bool   //popper配置 删除已出队完成的消息文件
	MsgHasSep        bool   //popper配置 是否返回消息结尾的符号
	ReadCount        bool   //popper配置 通过消息数量读取 否则通过偏移量读取
}

//获得消息文件名
func (m *Options) getMsgFileName(index int64) string {
	name := internal.Sprintf("%s.%d", m.FileNamePrefix, index)
	filename := internal.Path.ProgramDirJoin(m.MsgFileDir, name)
	return filename
}

//获得配置文件名
func (m *Options) getConfFileName(typ string) string {
	name := basal.Sprintf("%s_%s.fq", typ, m.FileNamePrefix)
	filename := basal.Path.ProgramDirJoin(m.ConfDataDir, name)
	return filename
}

//初始化未设置默认参数
func (m *Options) init() {
	if m.ConfDataDir == "" {
		m.ConfDataDir = "file_queue_conf"
	}
	if m.MsgFileDir == "" {
		m.MsgFileDir = "file_queue_msg"
	}
	if m.FileNamePrefix == "" {
		m.FileNamePrefix = "file_queue"
	}
	if len(m.Sep) == 0 {
		m.Sep = []byte("\n")
	}
	if m.MsgFileMaxByte < 1 {
		m.MsgFileMaxByte = MBToByteCount(30)
	}
	if m.PushChanSize < 1 {
		m.PushChanSize = 1000
	}
}

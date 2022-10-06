package filequeue

import (
	"github.com/jingyanbin/core/basal"
	internal "github.com/jingyanbin/core/internal"
	"time"
)

type Option struct {
	ConfDataDir       string        //配置文件目录
	Name              string        //名称
	MsgFileMaxByte    int64         //pusher配置 消息文件最大占用空间 单位: 字节
	PushChanSize      int           //pusher配置 入队管道大小
	DeletePoppedFile  bool          //popper配置 删除已出队完成的消息文件
	PrintInfoInterval time.Duration //打印信息间隔 0表示不打印
}

// 获得消息文件名
func (m *Option) getMsgFileName(index int64) string {
	name := internal.Sprintf("data.%d", index)
	filename := internal.Path.ProgramDirJoin(m.ConfDataDir, m.Name, "data", name)
	//filename := internal.Path.ProgramDirJoin(m.MsgFileDir, name)
	return filename
}

// 获得配置文件名
func (m *Option) getConfFileName(typ string) string {
	name := basal.Sprintf("%s.fq", typ)
	filename := internal.Path.ProgramDirJoin(m.ConfDataDir, m.Name, name)
	//filename := basal.Path.ProgramDirJoin(m.ConfDataDir, name)
	return filename
}

// 初始化未设置默认参数
func (m *Option) init() {
	if m.ConfDataDir == "" {
		m.ConfDataDir = "file_queue"
	}
	if m.Name == "" {
		m.Name = "queue1"
	}
	if m.MsgFileMaxByte < 1 {
		m.MsgFileMaxByte = MBToByteCount(30)
	}
	if m.PushChanSize < 1 {
		m.PushChanSize = 1000
	}
}

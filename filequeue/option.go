package filequeue

import (
	"github.com/jingyanbin/core/basal"
	internal "github.com/jingyanbin/core/internal"
	"path/filepath"
)

type Option struct {
	ConfDataDir      string //配置文件目录(配置非绝对路径都将自动转换为绝对路径)
	Name             string //名称(不要使用汉字)
	MsgFileMaxByte   int64  //pusher配置 消息文件最大占用空间 单位: 字节
	PushChanSize     int    //pusher配置 入队管道大小
	PushBufferSize   int    //pusher配置 缓冲大小
	DeletePoppedFile bool   //popper配置 删除已出队完成的消息文件
	PrintInfoSec     int    //多少秒打印一次信息 0表示不打印
}

// 获得消息文件名
func (m *Option) getMsgFileName(index int64) string {
	name := internal.Sprintf("data.%d", index)
	filename := internal.Path.Join(m.ConfDataDir, m.Name, "data", name)

	return filename
}

// 获得配置文件名
func (m *Option) getConfFileName(typ string) string {
	name := basal.Sprintf("%s.fq", typ)
	filename := internal.Path.Join(m.ConfDataDir, m.Name, name)
	return filename
}

// 初始化未设置默认参数
func (m *Option) init() {
	if m.ConfDataDir == "" {
		m.ConfDataDir = "file_queue"
	}
	if !filepath.IsAbs(m.ConfDataDir) {
		m.ConfDataDir = internal.Path.ProgramDirJoin(m.ConfDataDir)
	}
	if m.Name == "" {
		m.Name = "default"
	}
	if m.MsgFileMaxByte < 1 {
		m.MsgFileMaxByte = MBToByteCount(30)
	}
	if m.PushChanSize < 1 {
		m.PushChanSize = 1000
	}
	if m.PushBufferSize < 1 {
		m.PushBufferSize = 8000
	}
	if m.PrintInfoSec < 0 {
		m.PrintInfoSec = 0
	}
}

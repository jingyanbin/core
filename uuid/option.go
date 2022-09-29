package uuid

import (
	"github.com/jingyanbin/core/datetime"
	"github.com/jingyanbin/core/internal"
	tz "github.com/jingyanbin/core/timezone"
	"time"
)

var log = internal.GetStdoutLogger()

func SetLogger(logger internal.ILogger) {
	log = logger
}

const defaultEpoch = 1577808000000 //北京时间 2020-01-01 00:00:00.000
const backTime = time.Second

type Option struct {
	//基本配置
	Epoch        int64 //时间开始点: 毫秒
	IndexBits    int64 //每毫秒生产ID位数
	WorkerIdBits int64 //进程ID位数
	TimeBits     int64 //时间位数
	TimeLatest   bool  //使用最新时间
	WorkerId     int64 //进程ID
	BccSeed      uint8 //bcc校验种子
	Fast         bool  //快速模式,使用原子操作

	//计算的参数-------------
	//Shift
	workerIdShift int64
	timeShift     int64
	//max
	indexMax     int64
	workerIdMax  int64
	timeValueMax int64
	timeValueMin int64
	//date time
	dateTimeMax string
	dateTimeMin string
}

func (m *Option) init() {
	if m.IndexBits < 1 {
		panic(internal.NewError("uuid option IndexBits less: 1, %d", m.IndexBits))
	}
	if m.WorkerIdBits < 1 {
		panic(internal.NewError("uuid option WorkerIdBits less: 1, %d", m.WorkerIdBits))
	}
	if m.TimeBits < 1 {
		panic(internal.NewError("uuid option TimeBits less: 1, %d", m.TimeBits))
	}
	if totalBits := m.IndexBits + m.WorkerIdBits + m.TimeBits; totalBits > 63 {
		panic(internal.NewError("uuid option total bits more: %d/63", totalBits))
	}
	if m.Epoch < 1 {
		m.Epoch = defaultEpoch
	}
	m.workerIdShift = m.IndexBits
	m.timeShift = m.IndexBits + m.WorkerIdBits
	//
	m.indexMax = 1<<m.IndexBits - 1
	m.workerIdMax = 1<<m.WorkerIdBits - 1
	m.timeValueMax = 1<<m.TimeBits - 1
	m.timeValueMin = -1 << m.TimeBits

	if m.WorkerId > m.workerIdMax || m.WorkerId < 0 {
		panic(internal.NewError("uuid option WorkerId out of range: 0~%d, %d", m.workerIdMax, m.WorkerId))
	}
	m.dateTimeMin = datetime.UnixToYmdHMS((m.Epoch+m.timeValueMin)/1000, tz.Local())
	m.dateTimeMax = datetime.UnixToYmdHMS((m.Epoch+m.timeValueMax)/1000, tz.Local())
	nYear := m.timeValueMax / (3600 * 24 * 366 * 1000)
	//log.InfoF("uuid option time range: %v, %v, nYear: %v", m.dateTimeMin, m.dateTimeMax, nYear)
	now := datetime.UnixMs()
	if now-m.Epoch < m.timeValueMin {
		panic(internal.NewError("uuid option now time less than time min: %v, nYear: %v", m.dateTimeMin, nYear))
	}
	if now-m.Epoch > m.timeValueMax {
		panic(internal.NewError("uuid option now time more than time max: %v, nYear: %v", m.dateTimeMax, nYear))
	}
}

func (m *Option) info() string {
	nYear := m.timeValueMax / (3600 * 24 * 366 * 1000)
	return internal.Sprintf("uuid option index max: %d, worker id max: %d, time range: %s, %s, nYear: %d", m.indexMax, m.workerIdMax, m.dateTimeMin, m.dateTimeMax, nYear)
}

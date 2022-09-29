package uuid

import (
	"encoding/binary"
	"encoding/hex"
	datetime "github.com/jingyanbin/core/datetime"
	internal "github.com/jingyanbin/core/internal"
	tz "github.com/jingyanbin/core/timezone"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Generator struct {
	opt      Option
	lastTime int64
	index    int64
	mu       sync.Mutex
}

func NewGenerator(opt Option) *Generator {
	opt.init()
	generator := &Generator{lastTime: datetime.UnixMs(), opt: opt}
	return generator
}

func (m *Generator) Info() string {
	return m.opt.info()
}

func (m *Generator) getTimeIndex() (int64, int64) {
	m.mu.Lock()
	nowTime := datetime.UnixMs()
	m.index = (m.index + 1) & m.opt.indexMax
	if m.index == 0 {
		m.lastTime = m.lastTime + 1
	}
	for nowTime < m.lastTime {
		interval := time.Duration(m.lastTime-nowTime) * time.Millisecond
		if backTime < interval { //时间回拨超过1秒 一般
			log.ErrorF("uuid back in time: %s", datetime.UnixToYmdHMS(m.lastTime/1000, tz.Local()))
		}
		time.Sleep(interval)
		nowTime = datetime.UnixMs()
	}
	curTime, curIndex := m.lastTime, m.index
	m.mu.Unlock()
	return curTime, curIndex
}

func (m *Generator) getTimeIndexLatest() (int64, int64) {
	m.mu.Lock()
	nowTime := datetime.UnixMs()
	for nowTime < m.lastTime {
		interval := time.Duration(m.lastTime-nowTime) * time.Millisecond
		if backTime < interval { //时间回拨超过1秒 一般
			log.ErrorF("uuid back in time: %s", datetime.UnixToYmdHMS(m.lastTime/1000, tz.Local()))
		}
		time.Sleep(interval)
		nowTime = datetime.UnixMs()
	}
	if nowTime == m.lastTime {
		m.index = (m.index + 1) & m.opt.indexMax
		if m.index == 0 {
			m.lastTime += 1
		}
	} else {
		m.index = 0
		m.lastTime = nowTime
	}
	curTime, curIndex := m.lastTime, m.index
	m.mu.Unlock()
	return curTime, curIndex
}

func (m *Generator) UUID() int64 {
	var ms, index int64
	if m.opt.TimeLatest {
		ms, index = m.getTimeIndexLatest()
	} else {
		ms, index = m.getTimeIndex()
	}
	timeValue := ms - m.opt.Epoch
	if timeValue > m.opt.timeValueMax || timeValue < m.opt.timeValueMin {
		panic(internal.NewError("uuid timeValue out of range: %s~%s, %s", m.opt.dateTimeMin, m.opt.dateTimeMax, datetime.UnixToYmdHMS(ms/1000, tz.Local())))
	}
	return timeValue<<m.opt.timeShift | m.opt.WorkerId<<m.opt.workerIdShift | index
}

func (m *Generator) UUIDStr() string {
	return strconv.FormatInt(m.UUID(), 10)
}

func (m *Generator) UUIDHex() string {
	return m.ToHex(m.UUID())
}

func (m *Generator) UUIDHexEx(ex ...byte) string {
	dLen := 8 + len(ex)
	data := make([]byte, dLen+1)
	binary.BigEndian.PutUint64(data[0:], uint64(m.UUID()))
	for i, c := range ex {
		data[8+i] = c
	}
	data[dLen] = m.bcc(data, 0, dLen)

	uuid := strings.ToUpper(hex.EncodeToString(data))
	return uuid
}

func (m *Generator) DeUUID(uuid int64) (ms, workerId, index int64) {
	ms = (uuid >> m.opt.timeShift) + m.opt.Epoch
	workerId = (uuid >> m.opt.workerIdShift) & m.opt.workerIdMax
	index = uuid & m.opt.indexMax
	return
}

func (m *Generator) ToHex(uuid int64) string {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data[0:], uint64(uuid))
	return strings.ToUpper(hex.EncodeToString(data))
}

func (m *Generator) ToUUID(uuidHex string) (int64, error) {
	data, err := hex.DecodeString(uuidHex)
	if err != nil {
		return 0, err
	}
	if len(data) != 16 {
		return 0, internal.NewError("s len error: %v", len(data))
	}
	return int64(binary.BigEndian.Uint64(data)), nil
}

func (m *Generator) bcc(buf []byte, offset int, length int) byte {
	value := m.opt.BccSeed
	for i := offset; i < offset+length; i++ {
		value ^= buf[i]
	}
	return value
}

func (m *Generator) DeUUIDHexEx(uuidHexEx string) (ms, workerId, index int64, ex []byte, err error) {
	var data []byte
	data, err = hex.DecodeString(uuidHexEx)
	if err != nil {
		return
	}

	bcc := m.bcc(data, 0, len(data)-1)
	if bcc != data[len(data)-1] {
		err = internal.NewError("uuid hex ex decode bcc check failed")
		return
	}
	uuid := int64(binary.BigEndian.Uint64(data[0:]))
	ms, workerId, index = m.DeUUID(uuid)
	ex = data[8 : len(data)-1]
	return
}

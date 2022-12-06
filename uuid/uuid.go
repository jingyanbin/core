package uuid

import (
	"encoding/binary"
	"encoding/hex"
	xtime2 "github.com/jingyanbin/core/datetime"
	log2 "github.com/jingyanbin/core/log"
	"github.com/jingyanbin/log"
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
	generator := &Generator{lastTime: xtime2.UnixMs(), opt: opt}
	return generator
}

func (m *Generator) Info() string {
	return m.opt.info()
}

func (m *Generator) getTimeIndex() (int64, int64) {
	m.mu.Lock()
	nowTime := xtime2.UnixMs()
	m.index = (m.index + 1) & m.opt.indexMax
	if m.index == 0 {
		m.lastTime = m.lastTime + 1
	}
	for nowTime < m.lastTime {
		interval := time.Duration(m.lastTime-nowTime) * time.Millisecond
		if interval > backTime { //时间回拨超过1秒 一般
			log2.Error("uuid back in time: %s, sleep: %vms", xtime2.UnixToYmdHMS(m.lastTime/1000, xtime2.Local()), interval.Milliseconds())
		}
		time.Sleep(interval)
		nowTime = xtime2.UnixMs()
	}
	curTime, curIndex := m.lastTime, m.index
	m.mu.Unlock()
	return curTime, curIndex
}

func (m *Generator) getTimeIndexLatest() (int64, int64) {
	m.mu.Lock()
	nowTime := xtime2.UnixMs()
	for nowTime < m.lastTime {
		interval := time.Duration(m.lastTime-nowTime) * time.Millisecond
		if interval > backTime { //时间回拨超过1秒 一般
			log2.Error("uuid back in time: %s, sleep: %vms", xtime2.UnixToYmdHMS(m.lastTime/1000, xtime2.Local()), interval.Milliseconds())
		}
		time.Sleep(interval)
		nowTime = xtime2.UnixMs()
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
	var unixMs, index int64
	if m.opt.TimeLatest {
		unixMs, index = m.getTimeIndexLatest()
	} else {
		unixMs, index = m.getTimeIndex()
	}
	timeValue := unixMs - m.opt.Epoch
	if timeValue > m.opt.timeValueMax || timeValue < m.opt.timeValueMin {
		panic(log.NewError("uuid timeValue out of range: %s~%s, %s", m.opt.dateTimeMin, m.opt.dateTimeMax, xtime2.UnixToYmdHMS(unixMs/1000, xtime2.Local())))
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

func (m *Generator) DeUUIDStr(uuidStr string) (ms, workerId, index int64, err error) {
	var uuid int
	uuid, err = strconv.Atoi(uuidStr)
	if err != nil {
		return 0, 0, 0, err
	}
	ms, workerId, index = m.DeUUID(int64(uuid))
	return
}

func (m *Generator) DeUUIDHex(uuidHex string) (ms, workerId, index int64, err error) {
	var uuid int64
	uuid, err = m.ToUUID(uuidHex)
	if err != nil {
		return 0, 0, 0, err
	}
	ms, workerId, index = m.DeUUID(uuid)
	return
}

func (m *Generator) DeUUIDHexEx(uuidHexEx string) (ms, workerId, index int64, ex []byte, err error) {
	var data []byte
	data, err = hex.DecodeString(uuidHexEx)
	if err != nil {
		return
	}
	bcc := m.bcc(data, 0, len(data)-1)
	if bcc != data[len(data)-1] {
		err = log.NewError("uuid hex ex decode bcc check failed")
		return
	}
	uuid := int64(binary.BigEndian.Uint64(data[0:]))
	ms, workerId, index = m.DeUUID(uuid)
	ex = data[8 : len(data)-1]
	return
}

func (m *Generator) ToHex(uuid int64) (uuidHex string) {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data[0:], uint64(uuid))
	return strings.ToUpper(hex.EncodeToString(data))
}

func (m *Generator) ToUUID(uuidHex string) (uuid int64, err error) {
	var data []byte
	data, err = hex.DecodeString(uuidHex)
	if err != nil {
		return 0, err
	}
	if len(data) != 16 {
		return 0, log.NewError("uuidHex len error: %v", len(data))
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

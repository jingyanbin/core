package uuid

import (
	"encoding/binary"
	"encoding/hex"
	datetime "github.com/jingyanbin/core/datetime"
	internal "github.com/jingyanbin/core/internal"
	tz "github.com/jingyanbin/core/timezone"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type IGenerator interface {
	UUID() int64
	UUIDStr() string
	UUIDHex() string
	UUIDHexEx(ex ...byte) string
	DeUUID(uuid int64) (unixMs, workerId, index int64)
	DeUUIDStr(uuidStr string) (unixMs, workerId, index int64, err error)
	DeUUIDHex(uuidHex string) (unixMs, workerId, index int64, err error)
	DeUUIDHexEx(uuidHexEx string) (unixMs, workerId, index int64, ex []byte, err error)
	ToUUID(uuidHex string) (uuid int64, err error)
	ToHex(uuid int64) (uuidHex string)
	Info() string
}

func NewIGenerator(opt Option, fast bool) IGenerator {
	if fast {
		return NewFastGenerator(opt)
	} else {
		return NewGenerator(opt)
	}
}

type FastGenerator struct {
	opt  Option
	uuid int64
}

func NewFastGenerator(opt Option) *FastGenerator {
	opt.init()
	generator := &FastGenerator{opt: opt}
	generator.init()
	return generator
}

func (m *FastGenerator) init() {
	unixMs := datetime.UnixMs()
	timeValue := unixMs - m.opt.Epoch
	m.uuid = (timeValue << m.opt.timeShift) | (m.opt.WorkerId << m.opt.workerIdShift) | 0
}

func (m *FastGenerator) Info() string {
	return m.opt.info()
}

func (m *FastGenerator) genUUID() int64 {
	var uuid, unixMsOld, unixMsNow, index int64
	for {
		uuid = atomic.LoadInt64(&m.uuid)
		unixMsOld = (uuid >> m.opt.timeShift) + m.opt.Epoch
		unixMsNow = datetime.UnixMs()
		index = ((uuid & m.opt.indexMax) + 1) & m.opt.indexMax
		if index == 0 {
			unixMsOld = unixMsOld + 1
		}
		if unixMsNow < unixMsOld {
			interval := time.Duration(unixMsOld-unixMsNow) * time.Millisecond
			if interval > backTime { //时间回拨超过1秒 一般
				log.ErrorF("uuid back in time: %s, sleep: %vms", datetime.UnixToYmdHMS(unixMsOld/1000, tz.Local()), interval.Milliseconds())
			}
			time.Sleep(interval)
		} else {
			timeValue := unixMsOld - m.opt.Epoch
			if timeValue > m.opt.timeValueMax || timeValue < m.opt.timeValueMin {
				log.ErrorF("uuid timeValue out of range: %s~%s, %s", m.opt.dateTimeMin, m.opt.dateTimeMax, datetime.UnixToYmdHMS(unixMsNow/1000, tz.Local()))
				time.Sleep(time.Second)
				continue
			}
			uuidNow := timeValue<<m.opt.timeShift | m.opt.WorkerId<<m.opt.workerIdShift | index
			if atomic.CompareAndSwapInt64(&m.uuid, uuid, uuidNow) {
				return uuidNow
			} else {
				time.Sleep(time.Millisecond)
			}
		}
	}
}

func (m *FastGenerator) genUUIDLast() int64 {
	var uuid, unixMsOld, unixMsNow, index int64
	for {
		uuid = atomic.LoadInt64(&m.uuid)
		unixMsOld = (uuid >> m.opt.timeShift) + m.opt.Epoch
		unixMsNow = datetime.UnixMs()
		if unixMsNow < unixMsOld {
			interval := time.Duration(unixMsOld-unixMsNow) * time.Millisecond
			if interval > backTime { //时间回拨超过1秒 一般
				log.ErrorF("uuid back in time: %s, sleep: %vms", datetime.UnixToYmdHMS(unixMsOld/1000, tz.Local()), interval.Milliseconds())
			}
			time.Sleep(interval)
		} else if unixMsNow == unixMsOld { //时间一样
			index = ((uuid & m.opt.indexMax) + 1) & m.opt.indexMax
			timeValue := unixMsNow - m.opt.Epoch
			if timeValue > m.opt.timeValueMax || timeValue < m.opt.timeValueMin {
				log.ErrorF("uuid timeValue out of range: %s~%s, %s", m.opt.dateTimeMin, m.opt.dateTimeMax, datetime.UnixToYmdHMS(unixMsNow/1000, tz.Local()))
				time.Sleep(time.Second)
				continue
			}
			uuidNow := timeValue<<m.opt.timeShift | m.opt.WorkerId<<m.opt.workerIdShift | index
			if atomic.CompareAndSwapInt64(&m.uuid, uuid, uuidNow) {
				return uuidNow
			} else {
				time.Sleep(time.Millisecond)
			}
		} else {
			timeValue := unixMsNow - m.opt.Epoch
			uuidNow := timeValue<<m.opt.timeShift | m.opt.WorkerId<<m.opt.workerIdShift | 0
			if atomic.CompareAndSwapInt64(&m.uuid, uuid, uuidNow) {
				return uuidNow
			} else {
				time.Sleep(time.Millisecond)
			}
		}
	}
}

func (m *FastGenerator) UUID() int64 {
	if m.opt.TimeLatest {
		return m.genUUIDLast()
	} else {
		return m.genUUID()
	}
}

func (m *FastGenerator) UUIDStr() string {
	return strconv.FormatInt(m.UUID(), 10)
}

func (m *FastGenerator) UUIDHex() string {
	return m.ToHex(m.UUID())
}

func (m *FastGenerator) UUIDHexEx(ex ...byte) string {
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

func (m *FastGenerator) DeUUID(uuid int64) (unixMs, workerId, index int64) {
	unixMs = (uuid >> m.opt.timeShift) + m.opt.Epoch
	workerId = (uuid >> m.opt.workerIdShift) & m.opt.workerIdMax
	index = uuid & m.opt.indexMax
	return
}

func (m *FastGenerator) DeUUIDStr(uuidStr string) (unixMs, workerId, index int64, err error) {
	var uuid int
	uuid, err = strconv.Atoi(uuidStr)
	if err != nil {
		return 0, 0, 0, err
	}
	unixMs, workerId, index = m.DeUUID(int64(uuid))
	return
}

func (m *FastGenerator) DeUUIDHex(uuidHex string) (unixMs, workerId, index int64, err error) {
	var uuid int64
	uuid, err = m.ToUUID(uuidHex)
	if err != nil {
		return 0, 0, 0, err
	}
	unixMs, workerId, index = m.DeUUID(uuid)
	return
}

func (m *FastGenerator) DeUUIDHexEx(uuidHexEx string) (unixMs, workerId, index int64, ex []byte, err error) {
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
	unixMs, workerId, index = m.DeUUID(uuid)
	ex = data[8 : len(data)-1]
	return
}

func (m *FastGenerator) ToHex(uuid int64) (uuidHex string) {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data[0:], uint64(uuid))
	return strings.ToUpper(hex.EncodeToString(data))
}

func (m *FastGenerator) ToUUID(uuidHex string) (uuid int64, err error) {
	var data []byte
	data, err = hex.DecodeString(uuidHex)
	if err != nil {
		return 0, err
	}
	if len(data) != 16 {
		return 0, internal.NewError("uuidHex len error: %v", len(data))
	}
	return int64(binary.BigEndian.Uint64(data)), nil
}

func (m *FastGenerator) bcc(buf []byte, offset int, length int) byte {
	value := m.opt.BccSeed
	for i := offset; i < offset+length; i++ {
		value ^= buf[i]
	}
	return value
}

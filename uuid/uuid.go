package uuid

import (
	"encoding/binary"
	"encoding/hex"
	datetime "github.com/jingyanbin/core/datetime"
	internal "github.com/jingyanbin/core/internal"
	tz "github.com/jingyanbin/core/timezone"
	"sync"
	"time"
)

var log = internal.GetStdoutLogger()

func SetLogger(logger internal.ILogger) {
	log = logger
}

const defaultEpoch = 1577808000000 //北京时间 2020-01-01 00:00:00.000
const backTime = 1000 * time.Millisecond

type Uint64Generator struct {
	lastTime int64
	index    int64
	mu       sync.Mutex
	latest   bool

	epoch    int64
	workerId int64

	workerIdBits int64
	indexBits    int64
	//timeBits     int64
	unsigned bool

	workerIdMax int64
	indexMax    int64
	timeMax     int64
	timeMin     int64

	indexShift int64
	timeShift  int64

	dateTimeMax string
	dateTimeMin string
}

func (m *Uint64Generator) getTimeIndex() (int64, int64) {
	m.mu.Lock()
	nowTime := datetime.UnixMs()
	m.index = (m.index + 1) & m.indexMax
	if m.index == 0 {
		m.lastTime = m.lastTime + 1
	}
	for nowTime < m.lastTime {
		interval := time.Duration(m.lastTime-nowTime) * time.Millisecond
		if backTime < interval { //时间回拨超过1秒 一般
			log.ErrorF("uint64 uuid getTimeIndex back in time: %s", datetime.UnixToYmdHMS(m.lastTime/1000, tz.Local()))
		}
		time.Sleep(interval)
		nowTime = datetime.UnixMs()
	}
	curTime, curIndex := m.lastTime, m.index
	m.mu.Unlock()
	return curTime, curIndex
}

func (m *Uint64Generator) getTimeIndexNow() (int64, int64) {
	m.mu.Lock()
	nowTime := datetime.UnixMs()
	for nowTime < m.lastTime {
		interval := time.Duration(m.lastTime-nowTime) * time.Millisecond
		if backTime < interval { //时间回拨超过1秒 一般
			log.ErrorF("uint64 uuid getTimeIndexNow back in time: %s", datetime.UnixToYmdHMS(m.lastTime/1000, tz.Local()))
		}
		time.Sleep(interval)
		nowTime = datetime.UnixMs()
	}
	if nowTime == m.lastTime {
		m.index = (m.index + 1) & m.indexMax
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

func (m *Uint64Generator) UUID() uint64 {
	var ms, index int64
	if m.latest {
		ms, index = m.getTimeIndexNow()
	} else {
		ms, index = m.getTimeIndex()
	}
	value := ms - m.epoch
	if value > m.timeMax || value < m.timeMin {
		panic(internal.NewError("uuid time out of range: %s~%s, %s", m.dateTimeMin, m.dateTimeMax, datetime.UnixToYmdHMS(ms/1000, tz.Local())))
	}
	return uint64((value << m.timeShift) | (index << m.indexShift) | m.workerId)
}

func (m *Uint64Generator) DeUUID(uuid uint64) (ms int64, index, workerId int) {
	if m.unsigned {
		ms = int64(uuid>>uint64(m.timeShift)) + m.epoch
		index = int((uuid >> uint64(m.indexShift)) & uint64(m.indexMax))
		workerId = int(uuid & uint64(m.workerIdMax))
	} else {
		uid := int64(uuid)
		ms = uid>>m.timeShift + m.epoch
		index = int((uid >> m.indexShift) & m.indexMax)
		workerId = int(uid & m.workerIdMax)
	}
	return
}

func (m *Uint64Generator) Hex() string {
	return m.ToHex(m.UUID())
}

func (m *Uint64Generator) ToHex(uuid uint64) string {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data[0:], uuid)
	return hex.EncodeToString(data)
}

func (m *Uint64Generator) ToUUID(s string) (uint64, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(data), nil
}

// NewUUIDGenerator workerId: 服务进程唯一ID 范围: 0~（1<<workerIdBits-1）
//workerIdBits: workerId 占用位数(1~21 与indexBits共占22bits indexBits越大并发支持越大,同时支持的workerId越少)
//epoch: 开始时间戳 毫秒
//latest: 是否使用最新时间生成ID
//unsigned: 是否使用无符号时间, true时不可表达时间小于epoch的时间, false时可表达小于epoch的时间
func NewUUIDGenerator(workerId, workerIdBits, epoch int, latest bool, unsigned bool) *Uint64Generator {
	if workerIdBits < 0 || workerIdBits > 21 {
		panic(internal.NewError("new uint64 generator workerIdBits out of range: 1~21, %d", workerIdBits))
	}
	workerIdMax := 1<<workerIdBits - 1
	if workerId > workerIdMax || workerId < 0 {
		panic(internal.NewError("new uint64 generator workerId out of range: 0~%d, %d", workerIdMax, workerId))
	}
	if epoch < 1 {
		epoch = defaultEpoch
	}
	ms := datetime.UnixMs()
	var timeMax, timeMin int64
	if unsigned {
		timeMax = 1<<42 - 1
		timeMin = 0
	} else {
		timeMax = 1<<41 - 1
		timeMin = -1 << 41
	}

	if ms-int64(epoch) < timeMin {
		panic(internal.NewError("new uint64 generator now time less than time min value: %d", timeMin))
	}

	if ms-int64(epoch) > timeMax {
		panic(internal.NewError("new uint64 generator now time more than time max value: %d", timeMax))
	}

	generator := &Uint64Generator{lastTime: datetime.UnixMs(), latest: latest, unsigned: unsigned}
	generator.epoch = int64(epoch)
	generator.workerId = int64(workerId)

	generator.workerIdBits = int64(workerIdBits)
	generator.indexBits = 22 - int64(workerIdBits)

	generator.workerIdMax = int64(workerIdMax)
	generator.indexMax = 1<<generator.indexBits - 1

	generator.indexShift = generator.workerIdBits
	generator.timeShift = 22 //generator.workerIdBits + generator.indexBits

	generator.timeMax = timeMax
	generator.timeMin = timeMin

	generator.dateTimeMax = datetime.UnixToYmdHMS(generator.epoch+generator.timeMax, tz.Local())
	generator.dateTimeMin = datetime.UnixToYmdHMS(generator.epoch+generator.timeMin, tz.Local())

	return generator
}

const hexIndexBits = 20
const hexTimeBits = 43
const hexTimeSignBits = 1

const hexIndexMax = 1<<hexIndexBits - 1
const hexTimeMax = 1<<hexTimeBits - 1
const hexTimeMin = -1 << hexIndexBits

const hexTimeShift = hexIndexBits

var hexDateTimeMax = datetime.UnixToYmdHMS(defaultEpoch+hexTimeMax, tz.Local())
var hexDateTimeMin = datetime.UnixToYmdHMS(defaultEpoch+hexTimeMin, tz.Local())

type HexGenerator struct {
	index int64
	mu    sync.Mutex

	lastTime int64
	workerId int64
	seed     byte
	incr     uint8
	latest   bool
}

func (m *HexGenerator) getTimeIndex() (int64, int64) {
	m.mu.Lock()
	nowTime := datetime.UnixMs()
	m.index = (m.index + 1) & hexIndexMax
	if m.index == 0 {
		m.lastTime = m.lastTime + 1
	}
	for nowTime < m.lastTime {
		interval := time.Duration(m.lastTime-nowTime) * time.Millisecond
		if backTime < interval { //时间回拨超过1秒 一般
			log.ErrorF("uint64 uuid getTimeIndex back in time: %s", datetime.UnixToYmdHMS(m.lastTime/1000, tz.Local()))
		}
		time.Sleep(interval)
		nowTime = datetime.UnixMs()
	}
	curTime, curIndex := m.lastTime, m.index
	m.mu.Unlock()
	return curTime, curIndex
}

func (m *HexGenerator) getTimeIndexNow() (int64, int64) {
	m.mu.Lock()
	nowTime := datetime.UnixMs()
	for nowTime < m.lastTime {
		interval := time.Duration(m.lastTime-nowTime) * time.Millisecond
		if backTime < interval { //时间回拨超过1秒 一般
			log.ErrorF("uint64 uuid hex generator getTimeIndexNow back in time: %s", datetime.UnixToYmdHMS(m.lastTime/1000, tz.Local()))
		}
		time.Sleep(interval)
		nowTime = datetime.UnixMs()
	}
	if nowTime == m.lastTime {
		m.index = (m.index + 1) & hexIndexMax
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

func (m *HexGenerator) bcc(buf []byte, offset int, length int) byte {
	value := m.seed
	for i := offset; i < offset+length; i++ {
		value ^= buf[i]
	}
	return value
}

func (m *HexGenerator) UUID() string {
	var ms, index int64
	if m.latest {
		ms, index = m.getTimeIndexNow()
	} else {
		ms, index = m.getTimeIndex()
	}
	timeValue := ms - defaultEpoch
	if timeValue > hexTimeMax || timeValue < hexTimeMin {
		panic(internal.NewError("uuid time out of range: %s~%s, %s", hexDateTimeMin, hexDateTimeMax, datetime.UnixToYmdHMS(ms/1000, tz.Local())))
	}
	data := make([]byte, 12)
	head := uint64((timeValue << hexTimeShift) | index)
	binary.BigEndian.PutUint32(data[0:], uint32(head))
	binary.BigEndian.PutUint16(data[4:], uint16(head>>32))
	binary.BigEndian.PutUint16(data[6:], uint16(head>>48))
	binary.BigEndian.PutUint16(data[8:], uint16(m.workerId))
	data[10] = m.incr
	data[11] = m.bcc(data, 0, 11)
	uid := hex.EncodeToString(data)
	return uid
}

func (m *HexGenerator) DeUUID(uuid string) (ms int64, index int, workerId uint16, err error) {
	if len(uuid) != 24 {
		err = internal.NewError("decode uuid length no is 26")
		return
	}
	var data []byte
	data, err = hex.DecodeString(uuid)
	if err != nil {
		return
	}
	bcc := m.bcc(data, 0, 11)
	if bcc != data[11] {
		err = internal.NewError("decode uuid bcc check failed")
		return
	}
	h1 := int64(binary.BigEndian.Uint32(data[0:]))
	h2 := int64(binary.BigEndian.Uint16(data[4:]))
	h3 := int64(binary.BigEndian.Uint16(data[6:]))
	head := h1 | h2<<32 | h3<<48
	ms = head>>hexTimeShift + defaultEpoch
	index = int(head & hexIndexMax)
	workerId = binary.BigEndian.Uint16(data[8:])
	return
}

func (m *HexGenerator) UUIDExtra(extra uint16) string {
	var ms, index int64
	if m.latest {
		ms, index = m.getTimeIndexNow()
	} else {
		ms, index = m.getTimeIndex()
	}
	timeValue := ms - defaultEpoch
	if timeValue > hexTimeMax || timeValue < hexTimeMin {
		panic(internal.NewError("uuid extra time out of range: %s~%s, %s", hexDateTimeMin, hexDateTimeMax, datetime.UnixToYmdHMS(ms/1000, tz.Local())))
	}
	data := make([]byte, 14)
	head := uint64((timeValue << hexTimeShift) | index)
	binary.BigEndian.PutUint32(data[0:], uint32(head))
	binary.BigEndian.PutUint16(data[4:], uint16(head>>32))
	binary.BigEndian.PutUint16(data[6:], uint16(head>>48))
	binary.BigEndian.PutUint16(data[8:], uint16(m.workerId))
	data[10] = m.incr
	binary.BigEndian.PutUint16(data[11:], extra)
	data[13] = m.bcc(data, 0, 13)
	uid := hex.EncodeToString(data)
	return uid
}

func (m *HexGenerator) DeUUIDExtra(uuid string) (ms int64, index int, workerId uint16, extra uint16, err error) {
	if len(uuid) != 28 {
		err = internal.NewError("decode uuid extra length no is 26")
		return
	}

	var data []byte
	data, err = hex.DecodeString(uuid)
	if err != nil {
		return
	}
	bcc := m.bcc(data, 0, 13)
	if bcc != data[13] {
		err = internal.NewError("decode uuid extra bcc check failed")
		return
	}
	h1 := int64(binary.BigEndian.Uint32(data[0:]))
	h2 := int64(binary.BigEndian.Uint16(data[4:]))
	h3 := int64(binary.BigEndian.Uint16(data[6:]))
	head := h1 | h2<<32 | h3<<48
	ms = head>>hexTimeShift + defaultEpoch
	index = int(head & hexIndexMax)
	workerId = binary.BigEndian.Uint16(data[8:])
	extra = binary.BigEndian.Uint16(data[11:])
	return
}

// NewHexGenerator workerId: 服务进程唯一ID 范围: 0~65535
//seed: bcc计算初始值
//latest: 是否使用最新时间生成ID
//incr: 单次启动id+1(避免时间回调重复生成id,如果反复回调时间超过incr上限还是会出现ID重复生成)
func NewHexGenerator(workerId uint16, seed uint8, incr uint8, latest bool) *HexGenerator {
	generator := &HexGenerator{lastTime: datetime.UnixMs(), workerId: int64(workerId), seed: seed, incr: incr, latest: latest}
	return generator
}

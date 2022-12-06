package filequeue

import (
	"encoding/binary"
	"github.com/jingyanbin/core/basal"
	"github.com/jingyanbin/core/internal"
	"github.com/jingyanbin/core/log"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// 弹出处理函数
// data:数据
// popped:已弹出
type PopHandler func(data []byte) (popped bool, exit bool)

// 文件队列弹出器
type fileQueuePopper struct {
	conf       configDataPopper //配置数据
	f          *os.File         //队列数文件
	lastData   []byte
	lastOffset int64
	nextFile   bool //是否有下一个文件
	closed     int32
	wg         sync.WaitGroup
	mu         sync.Mutex
}

func (m *fileQueuePopper) Count() int64 {
	return m.conf.count
}

func (m *fileQueuePopper) CurOffset() (int64, error) {
	return m.f.Seek(0, io.SeekCurrent)
}

func (m *fileQueuePopper) read() (int64, []byte, error) {
	if m.nextFile {
		if !m.isExistNext() {
			return 0, nil, io.EOF //下一个文件不存在, 表示本文件读到EOF
		}
		if err := m.openNext(); err != nil {
			return 0, nil, err //打开下一个文件失败,表示本文件读到EOF
		}
	} else {
		if err := m.reopen(false); err != nil {
			return 0, nil, err
		}
	}
	offset, err := m.CurOffset()
	if err != nil {
		return 0, nil, err
	}
	head := make([]byte, 5)
	n, err := m.f.Read(head[:1])
	if err != nil && n != 1 {
		if _, errSeek := m.f.Seek(offset, io.SeekStart); errSeek != nil {
			err = errSeek
		}
		return 0, nil, err
	}
	if head[0] == msgEOF { //读到文件末尾
		m.nextFile = true
		if _, errSeek := m.f.Seek(offset, io.SeekStart); errSeek != nil {
			err = errSeek
		}
		return 0, nil, io.EOF
	}
	n, err = m.f.Read(head[1:])
	if err != nil && n != 4 {
		if _, errSeek := m.f.Seek(offset, io.SeekStart); errSeek != nil {
			err = errSeek
		}
		return 0, nil, err
	}
	dLen := int(binary.BigEndian.Uint32(head[1:]))
	data := make([]byte, dLen)
	n, err = io.ReadFull(m.f, data)
	if err != nil && n != dLen {
		if _, errSeek := m.f.Seek(offset, io.SeekStart); errSeek != nil {
			err = errSeek
		}
		return 0, nil, err
	}
	//offset, err := m.f.Seek(0, io.SeekCurrent)
	//if err != nil {
	//	return 0, nil, err
	//}
	//if offset != preOffset+int64(dLen)+5 {
	//	panic("offset error")
	//}
	offset += int64(dLen) + 5
	dLen = len(data)
	if dLen > 0 {
		if data[dLen-1] == '\n' {
			data = data[:dLen-1]
		}
	}
	return offset, data, nil
}

func (m *fileQueuePopper) deleteMsgFile(index int64) {
	if m.conf.option.DeletePoppedFile {
		filename := m.conf.option.getMsgFileName(index)
		os.Remove(filename)
	}
}

func (m *fileQueuePopper) closeFile() {
	if m.f != nil {
		m.f.Close()
		m.f = nil
	}
}

func (m *fileQueuePopper) isExistNext() bool {
	filename := m.conf.option.getMsgFileName(m.conf.index + 1)
	has, err := basal.IsExist(filename)
	if err != nil {
		log.Error("FileQueuePopper isExistNext error: %v, %v", err, filename)
		return false
	}
	return has
}

func (m *fileQueuePopper) openNext() error {
	index := m.conf.index         //当前文件的编号
	nextIndex := m.conf.index + 1 //下一个文件编号
	f, err := basal.OpenFileB(m.conf.option.getMsgFileName(nextIndex), os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	if err = m.conf.SaveEx(nextIndex, 0); err != nil { //保存读新文件编号和偏移量
		return err
	}
	m.closeFile()          //关闭当前读的文件
	m.deleteMsgFile(index) //删除当前读的文件
	return m.setFile(f)    //设置新文件
}

func (m *fileQueuePopper) reopen(force bool) error {
	if m.f == nil || force {
		filename := m.conf.option.getMsgFileName(m.conf.index)
		has, err := basal.IsExist(filename)
		if err != nil {
			log.Error("FileQueuePopper reopen error: %v", err)
		}
		if !has {
			return io.EOF
		}
		m.closeFile()
		f, err := basal.OpenFileB(m.conf.option.getMsgFileName(m.conf.index), os.O_RDONLY, 0666)
		if err != nil {
			return err
		}
		return m.setFile(f) //重新打开
	}
	return nil
}

func (m *fileQueuePopper) openReset() {
	m.nextFile = false
}

func (m *fileQueuePopper) setFile(f *os.File) error {
	if _, err := f.Seek(m.conf.offset, io.SeekStart); err != nil {
		return err
	}
	if m.f != nil {
		m.f.Close()
	}
	m.f = f
	m.nextFile = false
	return nil
}

// 当前偏移量 不是文件偏移量, 这里是成功出队后的偏移量
func (m *fileQueuePopper) Offset() int64 {
	return m.conf.offset
}

// 丢弃队头数据 必须调用过 Front 有数据才能丢弃
func (m *fileQueuePopper) DiscardFront() (bool, error) {
	if m.lastData != nil {
		m.conf.count += 1
		m.lastData = nil
		m.nextFile = false
		m.conf.offset = m.lastOffset
		err := m.conf.Save()
		return true, err
	}
	return false, nil
}

func (m *fileQueuePopper) PopFrontBlock() (line []byte, ok bool) {
	interval := time.Millisecond
	var err error
	for atomic.LoadInt32(&m.closed) == 0 {
		if line, err = m.PopFront(); err == nil {
			return line, true
		} else if err == io.EOF {
			if interval < time.Second {
				interval += time.Millisecond * 200
			}
		} else {
			interval = time.Second
			log.Error("FileQueuePopper PopFrontBlock error: %v", err)
		}
		if interval > 0 {
			time.Sleep(interval)
		}
		//log.ErrorF("FileQueuePopper PopFrontBlock error========%v", err)
	}
	return nil, false
}

// 直接弹出队头数据
func (m *fileQueuePopper) PopFront() (line []byte, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	msg, err := m.Front()
	if err != nil {
		return nil, err
	}
	ok, err := m.DiscardFront()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, internal.NewError("DiscardFront failed")
	}
	return msg, nil
}

type Message struct {
	Offset int64  //偏移量
	Data   []byte //数据
}

// 获得队头数据
func (m *fileQueuePopper) Front() ([]byte, error) {
	if m.lastData != nil {
		return m.lastData, nil
	}
	offset, data, err := m.read()
	if err != nil {
		return nil, err
	}
	m.lastData = data
	m.lastOffset = offset
	return data, nil
}

func (m *fileQueuePopper) PopToHandler(handler PopHandler) {
	m.wg.Add(1)
	go m.popTo(handler)
}

func (m *fileQueuePopper) exit() {
	defer m.wg.Done()
	m.conf.Sync()
	m.conf.Close()
	if m.f != nil {
		m.f.Close()
	}
}

func (m *fileQueuePopper) popTo(handler PopHandler) {
	var interval time.Duration
	defer m.exit()
	var exit bool
	for atomic.LoadInt32(&m.closed) == 0 {
		m.mu.Lock()
		if data, err := m.Front(); err == nil {
			popped := false
			internal.Try(func() {
				popped, exit = handler(data)
			}, func(stack string, e error) {
				popped = false
				exit = false
				log.Error("FileQueuePopper popTo error: %v, %v", stack, string(data))
			})
			if exit {
				atomic.StoreInt32(&m.closed, 1)
			}
			if popped {
				ok, err2 := m.DiscardFront()
				if !ok || err != nil {
					log.Error("FileQueuePopper popTo DiscardFront error: %v, %v", ok, err2)
				}
				interval = 0
			} else {
				interval = time.Second
			}
		} else if err == io.EOF {
			if interval < time.Second {
				interval += time.Millisecond * 200
			}
		} else {
			interval = time.Second
			log.Error("FileQueuePopper popTo error: %v", err)
		}
		m.mu.Unlock()
		if interval > 0 {
			time.Sleep(interval)
		}
	}
}

func (m *fileQueuePopper) Closed() bool {
	return atomic.LoadInt32(&m.closed) == 1
}

func (m *fileQueuePopper) Close() {
	atomic.CompareAndSwapInt32(&m.closed, 0, 1)
	m.Wait()
	return
}

func (m *fileQueuePopper) Wait() {
	m.wg.Wait()
}

func newFileQueuePopper(option *Option) (*fileQueuePopper, error) {
	popper := &fileQueuePopper{}
	popper.conf.option = option
	popper.conf.filename = option.getConfFileName("pop")
	err := popper.conf.Load()
	if err != nil {
		return nil, err
	}
	return popper, nil
}

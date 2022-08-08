package filequeue

import (
	"github.com/jingyanbin/core/basal"
	internal "github.com/jingyanbin/core/internal"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const popperBufSize = 4096

//弹出处理函数
//data:数据
//popped:已弹出
type PopHandler func(data []byte) (popped bool, exit bool)

func (m *FileQueuePopper) splitHandlerByte(lines []byte, r int, sep byte) (pos int, nextStart int, canPop bool) {
	switch lines[r] {
	case '\n':
		if r > 0 && lines[r-1] == '\r' {
			pos = r - 1
		} else {
			pos = r
		}

		return pos, r + 1, true
	case sep:
		pos = r
		return pos, r + 1, true
	default:
		return 0, 0, false
	}
}

func (m *FileQueuePopper) splitHandlerBytes(lines []byte, r int, sep []byte) (pos int, nextStart int, canPop bool) {
	sepLen := len(sep)
	sepIndexMax := sepLen - 1
	for i := 0; i < sepLen; i++ {
		if lines[r-i] != sep[sepIndexMax-i] {
			return 0, 0, false
		}
	}
	pos = r - sepIndexMax
	nextStart = r + 1
	return pos, nextStart, true
}

//lines 已读出数据
//r 读出的位置
func (m *FileQueuePopper) splitHandler(lines []byte, r int, sep []byte) (pos int, nextStart int, canPop bool, nextFile bool) {
	if lines[r] == msgEOF { //数据末尾
		return 0, 0, false, true
	}
	sepLen := len(sep)
	if sepLen == 1 {
		pos, nextStart, canPop = m.splitHandlerByte(lines, r, sep[0])
		return pos, nextStart, canPop, false
	} else {
		linesLen := len(lines)
		if linesLen < sepLen {
			return 0, 0, false, false
		}
		pos, nextStart, canPop = m.splitHandlerBytes(lines, r, sep)
		return pos, nextStart, canPop, false
	}
}

//文件队列弹出器
type FileQueuePopper struct {
	conf      configDataPopper //配置数据
	f         *os.File         //队列数文件
	buf       []byte           //临时读取缓冲
	lines     []byte           //已读出数据 未出队
	r         int              //当前已读出下标
	pos       int              //待出队数据结束位置
	nextStart int              //下一次开始位置
	canPop    bool             //是否已经有待出队数据
	nextFile  bool             //是否有下一个文件
	readCount int64            //pop 当前读取数量
	closed    int32
	wg        sync.WaitGroup
	mu        sync.Mutex
}

func (m *FileQueuePopper) read() (n int, err error) {
	if m.nextFile {
		if !m.isExistNext() {
			return 0, io.EOF //下一个文件不存在, 表示本文件读到EOF
		}
		if err = m.openNext(); err != nil {
			return 0, err //打开下一个文件失败,表示本文件读到EOF
		}
	} else {
		if err = m.reopen(false); err != nil {
			return 0, err
		}
	}
	n, err = m.f.Read(m.buf)
	if err != nil && n == 0 {
		if err == io.EOF {
			return 0, err
		}
		if err = m.reopen(true); err != nil {
			return 0, err
		}
		n, err = m.f.Read(m.buf)
	}
	if n > 0 {
		m.lines = append(m.lines, m.buf[:n]...)
	}
	return
}

func (m *FileQueuePopper) deleteMsgFile(index int64) {
	if m.conf.options.DeletePoppedFile {
		filename := m.conf.options.getMsgFileName(index)
		os.Remove(filename)
	}
}

func (m *FileQueuePopper) closeFile() {
	if m.f != nil {
		m.f.Close()
		m.f = nil
	}
}

func (m *FileQueuePopper) isExistNext() bool {
	filename := m.conf.options.getMsgFileName(m.conf.index + 1)
	has, err := internal.IsExist(filename)
	if err != nil {
		log.ErrorF("FileQueuePopper isExistNext error: %v, %v", err, filename)
		return false
	}
	return has
}

func (m *FileQueuePopper) openNext() error {
	index := m.conf.index         //当前文件的编号
	nextIndex := m.conf.index + 1 //下一个文件编号
	f, err := basal.OpenFileB(m.conf.options.getMsgFileName(nextIndex), os.O_RDONLY, 0666)
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

func (m *FileQueuePopper) reopen(force bool) error {
	if m.f == nil || force {
		filename := m.conf.options.getMsgFileName(m.conf.index)
		has, err := basal.IsExist(filename)
		if err != nil {
			log.ErrorF("FileQueuePopper reopen error: %v", err)
		}
		if !has {
			return io.EOF
		}
		m.closeFile()
		f, err := basal.OpenFileB(m.conf.options.getMsgFileName(m.conf.index), os.O_RDONLY, 0666)
		if err != nil {
			return err
		}
		return m.setFile(f) //重新打开
	}
	return nil
}

func (m *FileQueuePopper) openReset() {
	m.lines = m.lines[:0]
	m.r = 0
	m.pos = 0
	m.nextStart = 0
	m.canPop = false
	m.nextFile = false
	m.readCount = 0
}

func (m *FileQueuePopper) setFile(f *os.File) error {
	if m.conf.options.ReadCount {
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return err
		}
	} else {
		if _, err := f.Seek(m.conf.offset, io.SeekStart); err != nil {
			return err
		}
	}
	if m.f != nil {
		m.f.Close()
	}
	m.f = f
	m.openReset()
	return nil
}

//当前偏移量 不是文件偏移量, 这里是成功出队后的偏移量
func (m *FileQueuePopper) Offset() int64 {
	return m.conf.offset
}

func (m *FileQueuePopper) discardReset() {
	m.lines = m.lines[m.nextStart:]
	m.r = 0
	m.pos = 0
	m.nextStart = 0
	m.canPop = false
	m.nextFile = false
}

//丢弃队头数据 必须调用过 Front 有数据才能丢弃
func (m *FileQueuePopper) DiscardFront() bool {
	if m.canPop {
		if m.conf.options.ReadCount {
			m.conf.offset = m.readCount
		} else {
			m.conf.offset += int64(m.nextStart)
		}
		m.discardReset()
		m.conf.Save()
		//log.ErrorF("================%v", m.readCount)
		return true
	} else {
		return false
	}
}

func (m *FileQueuePopper) PopFrontBlock() (line []byte, ok bool) {
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
			log.ErrorF("FileQueuePopper PopFrontBlock error: %v", err)
		}
		if interval > 0 {
			time.Sleep(interval)
		}
	}
	return nil, false
}

//直接弹出队头数据
func (m *FileQueuePopper) PopFront() (line []byte, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if line, err = m.Front(); err == nil && m.canPop {
		m.DiscardFront()
	}
	return
}

//获得队头数据
func (m *FileQueuePopper) Front() (line []byte, err error) {
	if m.canPop {
		return m.lines[:m.pos], nil
	}
	var n int
	for {
		for ; (m.r < len(m.lines)) && (m.nextFile == false); m.r++ {
			m.pos, m.nextStart, m.canPop, m.nextFile = m.splitHandler(m.lines, m.r, m.conf.options.Sep)
			if m.canPop {
				if m.conf.options.MsgHasSep {
					m.pos = m.nextStart
				}
				if m.conf.options.ReadCount {
					m.readCount += 1
					if m.readCount <= m.conf.offset {
						m.discardReset()
						continue
					}
				}
				return m.lines[:m.pos], nil
			}
		}
		n, err = m.read()
		if n == 0 {
			return nil, err
		}
	}
}

func (m *FileQueuePopper) PopToHandler(handler PopHandler) {
	m.wg.Add(1)
	go m.popTo(handler)
}

func (m *FileQueuePopper) exit() {
	defer m.wg.Done()
	m.conf.Sync()
	m.conf.Close()
	if m.f != nil {
		m.f.Close()
	}
}

func (m *FileQueuePopper) popTo(handler PopHandler) {
	var interval time.Duration
	defer m.exit()
	var exit bool
	for atomic.LoadInt32(&m.closed) == 0 {
		m.mu.Lock()
		if data, err := m.Front(); err == nil && m.canPop {
			popped := false
			basal.Try(func() {
				popped, exit = handler(data)
			}, func(stack string, e error) {
				popped = false
				exit = false
				log.ErrorF("FileQueuePopper popTo error: %v, %v", stack, string(data))
			})
			if exit {
				atomic.StoreInt32(&m.closed, 1)
			}
			if popped {
				m.DiscardFront()
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
			log.ErrorF("FileQueuePopper popTo error: %v", err)
		}
		m.mu.Unlock()
		if interval > 0 {
			time.Sleep(interval)
		}
	}
}

func (m *FileQueuePopper) Close() {
	atomic.CompareAndSwapInt32(&m.closed, 0, 1)
	m.Wait()
	return
}

func (m *FileQueuePopper) Wait() {
	m.wg.Wait()
}

func NewFileQueuePopper(options *Options) (*FileQueuePopper, error) {
	if options == nil {
		options = &Options{}
	}
	options.init()
	popper := &FileQueuePopper{}
	popper.conf.options = options
	popper.conf.filename = options.getConfFileName("pop")
	err := popper.conf.Load()
	if err != nil {
		return nil, err
	}
	popper.buf = make([]byte, popperBufSize)
	popper.lines = make([]byte, 0, popperBufSize)
	return popper, nil
}

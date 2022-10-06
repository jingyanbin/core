package filequeue

import (
	"encoding/binary"
	"github.com/jingyanbin/core/basal"
	internal "github.com/jingyanbin/core/internal"
	"os"
	"sync"
	"sync/atomic"
)

// FileQueuePusher 文件队列入队器
type FileQueuePusher struct {
	conf   configDataPusher //配置数据
	ch     chan []byte      //缓冲
	closed int32            //关闭状态
	wg     sync.WaitGroup   //关闭等待组
	f      *os.File         //当前写入文件
	count  int64            //当前写入条数
}

func (m *FileQueuePusher) ChanLenAndSize() (int, int) {
	return len(m.ch), m.conf.option.PushChanSize
}

func (m *FileQueuePusher) Count() int64 {
	return m.count
}

// 当前大小
func (m *FileQueuePusher) size() (size int64, err error) {
	if err = m.reopen(false); err != nil {
		return 0, err
	}
	var fi os.FileInfo
	if fi, err = m.f.Stat(); err != nil {
		if err = m.reopen(true); err != nil {
			return 0, err
		}
		if fi, err = m.f.Stat(); err != nil {
			return 0, err
		}
	}
	return fi.Size(), nil
}

func (m *FileQueuePusher) write(buf []byte) (n int, err error) {
	size := len(buf)
	var nn int
	for n < size && err == nil {
		nn, err = m.f.Write(buf[n:])
		n += nn
	}
	return
}

func (m *FileQueuePusher) reopen(force bool) error {
	if m.f == nil || force {
		if m.f != nil {
			m.f.Sync()
			m.f.Close()
		}
		f, err := internal.OpenFileB(m.conf.option.getMsgFileName(m.conf.index), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			m.f = nil
			return err
		}
		m.f = f
	}
	return nil
}

func (m *FileQueuePusher) push(data []byte) (err error) {
	var size int64
	if size, err = m.size(); err != nil {
		return err
	}
	if m.conf.option.MsgFileMaxByte > 0 && size > m.conf.option.MsgFileMaxByte {
		if _, err = m.write([]byte{msgEOF}); err != nil { //类型为27是文件结束符
			return err
		}
		if err = m.conf.Next(); err != nil {
			return err
		}
		if err = m.reopen(true); err != nil {
			return err
		}
	}
	dLen := len(data)
	if dLen == 0 {
		return basal.NewError("数据为空")
	}
	//if data[dLen-1] != '\n' {
	//	data = append(data, '\n')
	//}
	data = append(data, '\n')
	dLen = len(data)
	pushData := make([]byte, 5, dLen+5)
	//pushData = append(pushData, 0) //类型为0是一般消息
	//pushData = binary.BigEndian.Uint32(pushData, uint32(dLen))
	binary.BigEndian.PutUint32(pushData[1:], uint32(dLen))
	pushData = append(pushData, data...)
	var n int
	if n, err = m.write(pushData); err != nil {
		if err = m.reopen(true); err != nil {
			return err
		}
		_, err = m.write(pushData[n:])
	}
	return err
}

func (m *FileQueuePusher) exit() {
	defer m.wg.Done()
	m.conf.Sync()
	m.conf.Close()
	if m.f != nil {
		m.f.Sync()
		m.f.Close()
	}
}

func (m *FileQueuePusher) run() {
	defer m.exit()
	var err error
	for data := range m.ch {
		err = m.push(data)
		if err != nil {
			log.ErrorF("FileQueuePusher run push error: %v, data: %v", err, string(data))
		}
		m.count += 1
	}
}

func (m *FileQueuePusher) Closed() bool {
	return atomic.LoadInt32(&m.closed) == 1
}

func (m *FileQueuePusher) Close() {
	if atomic.CompareAndSwapInt32(&m.closed, 0, 1) {
		close(m.ch)
	}
	m.Wait()
	return
}

func (m *FileQueuePusher) Wait() {
	m.wg.Wait()
}

func (m *FileQueuePusher) Push(data []byte) (err error) {
	defer basal.Exception(func(stack string, e error) {
		err = e
	})
	//if pushLen := len(m.ch); pushLen >= m.conf.options.PushChanSize-5 {
	//	log.ErrorF("FileQueue Push full: %d", pushLen)
	//}
	m.ch <- data
	return
}

func (m *FileQueuePusher) PushString(data string) error {
	return m.Push([]byte(data))
}

func newFileQueuePusher(option Option) (*FileQueuePusher, error) {
	pusher := &FileQueuePusher{}
	pusher.conf.option = option
	pusher.conf.filename = option.getConfFileName("push")
	err := pusher.conf.Load()
	if err != nil {
		return nil, err
	}
	pusher.ch = make(chan []byte, option.PushChanSize)
	pusher.wg.Add(1)
	go pusher.run()
	return pusher, nil
}

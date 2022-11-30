package filequeue

import (
	"encoding/binary"
	"github.com/jingyanbin/core/basal"
	internal "github.com/jingyanbin/core/internal"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var errDataNil = basal.NewError("数据为空")

type batchBuffer struct {
	buf   []byte
	count int64
}

func newBatchBuffer(size int) *batchBuffer {
	return &batchBuffer{buf: make([]byte, 0, size)}
}

func (m *batchBuffer) Count() int64 {
	return m.count
}

func (m *batchBuffer) Len() int {
	return len(m.buf)
}

func (m *batchBuffer) Bytes() []byte {
	return (*m).buf
}

func (m *batchBuffer) Add(data []byte) error {
	dLen := len(data)
	if dLen == 0 {
		return errDataNil
	}
	dLen += 1                        //1字节结束符\n
	pushData := make([]byte, 5+dLen) //1字节数据类型,4字节uint32数据长度,数据
	binary.BigEndian.PutUint32(pushData[1:], uint32(dLen))
	copy(pushData[5:], data)
	pushData[4+dLen] = '\n'
	(*m).buf = append((*m).buf, pushData...)
	(*m).count += 1
	return nil
}

func (m *batchBuffer) Clear() {
	(*m).buf = (*m).buf[:0]
	(*m).count = 0
}

// 文件队列入队器
type fileQueuePusher struct {
	conf   configDataPusher //配置数据
	ch     chan []byte      //缓冲
	closed int32            //关闭状态
	wg     sync.WaitGroup   //关闭等待组
	f      *os.File         //当前写入文件
}

func (m *fileQueuePusher) ChanLenAndSize() (int, int) {
	return len(m.ch), m.conf.option.PushChanSize
}

func (m *fileQueuePusher) Count() int64 {
	return m.conf.count
}

// 当前大小
func (m *fileQueuePusher) size() (size int64, err error) {
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
	} else {
		if !internal.IsExistByFileInfo(fi) {
			if err = m.reopen(true); err != nil {
				return 0, err
			}
			if fi, err = m.f.Stat(); err != nil {
				return 0, err
			}
		}
	}
	return fi.Size(), nil
}

func (m *fileQueuePusher) write(buf []byte) (n int, err error) {
	size := len(buf)
	var nn int
	for n < size && err == nil {
		nn, err = m.f.Write(buf[n:])
		n += nn
	}
	return
}

func (m *fileQueuePusher) reopen(force bool) error {
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

func (m *fileQueuePusher) push(data []byte) (err error) {
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
		return errDataNil
	}
	var n int
	if n, err = m.write(data); err != nil {
		if err = m.reopen(true); err != nil {
			return err
		}
		_, err = m.write(data[n:])
	}
	return err
}

func (m *fileQueuePusher) pushOne(data []byte) (err error) {
	dLen := len(data)
	if dLen == 0 {
		return errDataNil
	}
	dLen += 1                        //1字节结束符\n
	pushData := make([]byte, 5+dLen) //1字节数据类型,4字节uint32数据长度,数据
	binary.BigEndian.PutUint32(pushData[1:], uint32(dLen))
	copy(pushData[5:], data)
	pushData[4+dLen] = '\n'
	return m.push(pushData)
}

func (m *fileQueuePusher) pushBatch(buf *batchBuffer) {
	if err := m.push(buf.Bytes()); err != nil {
		internal.Log.Error("FileQueuePusher pushBuffer error: %v, data: %v", err, string(buf.Bytes()))
	} else {
		m.conf.AddCount(buf.Count())
	}
	buf.Clear()
}

func (m *fileQueuePusher) exit() {
	defer m.wg.Done()
	m.conf.Sync()
	m.conf.Close()
	if m.f != nil {
		m.f.Sync()
		m.f.Close()
	}
}

func (m *fileQueuePusher) run() {
	defer m.exit()
	var err error
	//批量
	bufSize := m.conf.option.PushBufferSize + (m.conf.option.PushBufferSize / 10 * 2)
	buf := newBatchBuffer(bufSize)
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	for {
		select {
		case data, ok := <-m.ch:
			if !ok {
				if buf.Len() > 0 { //退出前发生剩余数据
					m.pushBatch(buf)
				}
				return
			}
			if err = buf.Add(data); err == nil {
				if buf.Len() > m.conf.option.PushBufferSize && buf.Count() > 10 { //长度超过限制发生
					m.pushBatch(buf)
				} else { //添加数据后没有push,就设置下次超时
					timer.Reset(time.Second)
				}
			} else {
				internal.Log.Error("FileQueuePusher run Add error: %v, data: %v", err, string(data))
			}

		case <-timer.C:
			if buf.Len() > 0 { //超时发送数据
				m.pushBatch(buf)
			}
		}
	}

	//单个
	//for data := range m.ch {
	//	err = m.pushOne(data)
	//	if err != nil {
	//		log.ErrorF("FileQueuePusher run push error: %v, data: %v", err, string(data))
	//	} else {
	//		m.conf.AddCount(1)
	//	}
	//}
}

func (m *fileQueuePusher) Closed() bool {
	return atomic.LoadInt32(&m.closed) == 1
}

func (m *fileQueuePusher) Close() {
	if atomic.CompareAndSwapInt32(&m.closed, 0, 1) {
		close(m.ch)
	}
	m.Wait()
	return
}

func (m *fileQueuePusher) Wait() {
	m.wg.Wait()
}

func (m *fileQueuePusher) Push(data []byte) (err error) {
	defer basal.Exception(func(stack string, e error) {
		err = e
	})
	//if pushLen := len(m.ch); pushLen >= m.conf.options.PushChanSize-5 {
	//	log.ErrorF("FileQueue Push full: %d", pushLen)
	//}
	m.ch <- data
	return
}

func (m *fileQueuePusher) PushString(data string) error {
	return m.Push([]byte(data))
}

func newFileQueuePusher(option *Option) (*fileQueuePusher, error) {
	pusher := &fileQueuePusher{}
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

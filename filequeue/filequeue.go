package filequeue

import (
	"github.com/jingyanbin/core/internal"
	"time"
)

type FileQueue struct {
	option *Option
	pusher *fileQueuePusher
	popper *fileQueuePopper
}

func (m *FileQueue) Info() string {
	name := internal.Path.Join(m.option.ConfDataDir, m.option.Name)
	chLen, chSize := m.pusher.ChanLenAndSize()
	c1, c2 := m.popper.Count(), m.pusher.Count()
	return internal.Sprintf("name: %s, push chan: %d/%d, popped/pushed: %d/%d, len: %d", name, chLen, chSize, c1, c2, c2-c1)
}

func (m *FileQueue) ClosePusher() {
	m.pusher.Close()
}

func (m *FileQueue) ClosePopper() {
	m.popper.Close()
}

func (m *FileQueue) Close() {
	m.popper.Close()
	m.pusher.Close()
	if m.pusher.conf.option.PrintInfoSec > 0 {
		internal.Log.Info("file queue close info: %s", m.Info())
	}
}

//func (m *FileQueue) Wait() {
//	m.pusher.Wait()
//	m.popper.Wait()
//}

func (m *FileQueue) Push(data []byte) error {
	return m.pusher.Push(data)
}

func (m *FileQueue) PopToHandler(handler PopHandler) {
	m.popper.PopToHandler(handler)
}

func (m *FileQueue) Pop() (data []byte, ok bool) {
	return m.popper.PopFrontBlock()
}

func (m *FileQueue) run(nSec int) {
	ticker := time.NewTicker(time.Duration(nSec) * time.Second)
	defer ticker.Stop()
	for !m.pusher.Closed() || !m.popper.Closed() {
		select {
		case <-ticker.C:
			internal.Log.Info("file queue run info: %s", m.Info())
		}
	}
}

// 创建文件消息队列
func NewFileQueue(option Option, popHandler PopHandler) (*FileQueue, error) {
	option.init()
	q := &FileQueue{option: &option}
	if pusher, err := newFileQueuePusher(q.option); err != nil {
		return nil, err
	} else {
		q.pusher = pusher
	}
	if popper, err := newFileQueuePopper(q.option); err != nil {
		return nil, err
	} else {
		q.popper = popper
	}
	if popHandler != nil {
		q.popper.PopToHandler(popHandler)
	}
	if option.PrintInfoSec > 0 {
		go q.run(option.PrintInfoSec)
		internal.Log.Info("file queue new info: %s", q.Info())
	}
	return q, nil
}

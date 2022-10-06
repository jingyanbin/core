package filequeue

import (
	"github.com/jingyanbin/core/internal"
	"time"
)

type FileQueue struct {
	pusher *FileQueuePusher
	popper *FileQueuePopper
}

func (m *FileQueue) Info() string {
	name := internal.Path.ProgramDirJoin(m.pusher.conf.option.ConfDataDir, m.pusher.conf.option.Name)
	chLen, chSize := m.pusher.ChanLenAndSize()
	return internal.Sprintf("file queue: %s, pushed/popped: %d/%d, push chan: %d/%d", name, m.pusher.Count(), m.popper.Count(), chLen, chSize)
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

func (m *FileQueue) run(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for !m.pusher.Closed() || !m.popper.Closed() {
		select {
		case <-ticker.C:
			log.InfoF(m.Info())
		}
	}
}

//创建文件消息队列
func NewFileQueue(option Option, popHandler PopHandler) (*FileQueue, error) {
	option.init()
	q := &FileQueue{}
	if pusher, err := newFileQueuePusher(option); err != nil {
		return nil, err
	} else {
		q.pusher = pusher
	}
	if popper, err := newFileQueuePopper(option); err != nil {
		return nil, err
	} else {
		q.popper = popper
	}
	if option.PrintInfoInterval > 0 {
		go q.run(option.PrintInfoInterval)
	}
	if popHandler != nil {
		q.popper.PopToHandler(popHandler)
	}
	log.InfoF("NewFileQueue info: %s", q.Info())
	return q, nil
}

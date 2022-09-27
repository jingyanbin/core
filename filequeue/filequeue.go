package filequeue

import "github.com/jingyanbin/core/internal"

type FileQueue struct {
	pusher *FileQueuePusher
	popper *FileQueuePopper
}

func (m *FileQueue) ClosePusher() {
	m.pusher.Close()
}

func (m *FileQueue) ClosePopper() {
	m.popper.Close()
}

func (m *FileQueue) Close() {
	m.pusher.Close()
	m.popper.Close()
}

func (m *FileQueue) Wait() {
	m.pusher.Wait()
	m.popper.Wait()
}

func (m *FileQueue) Push(data []byte) error {
	return m.pusher.Push(data)
}

func (m *FileQueue) PopToHandler(handler PopHandler) {
	m.popper.PopToHandler(handler)
}

func (m *FileQueue) Pop() (data []byte, ok bool) {
	return m.popper.PopFrontBlock()
}

//创建文件消息队列
func NewFileQueue(options *Options) (*FileQueue, error) {
	q := &FileQueue{}
	if pusher, err := NewFileQueuePusher(options); err != nil {
		return nil, err
	} else {
		q.pusher = pusher
	}
	if popper, err := NewFileQueuePopper(options); err != nil {
		return nil, err
	} else {
		q.popper = popper
	}
	log.InfoF("NewFileQueue %s", internal.Path.ProgramDirJoin(options.ConfDataDir, options.Name))
	return q, nil
}

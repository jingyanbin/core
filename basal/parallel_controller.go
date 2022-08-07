package basal

// ParallelController 并发控制器
type ParallelController struct {
	ch chan struct{}
}

// Acquire 获取
func (m *ParallelController) Acquire(block bool) bool {
	if block {
		m.ch <- struct{}{}
		return true
	} else {
		select {
		case m.ch <- struct{}{}:
			return true
		default:
			return false
		}
	}
}

// Release 释放
func (m *ParallelController) Release(block bool) bool {
	if block {
		<-m.ch
		return true
	} else {
		select {
		case <-m.ch:
			return true
		default:
			return false
		}
	}
}

func NewParallelController(max int) *ParallelController {
	return &ParallelController{ch: make(chan struct{}, max)}
}

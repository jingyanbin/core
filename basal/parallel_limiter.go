package basal

import (
	"sync/atomic"
	"time"
	"unsafe"
)

// ParallelLimiter
// @Description: 并发限制器
type ParallelLimiter struct {
	c *ParallelController
}

// Init
//
//	@Description: 初始化
//	@receiver m
//	@param max 最大并发数
//	@param block 是否苏塞
//	@param timeout 超时时间
func (m *ParallelLimiter) Init(max int32, block bool, timeout time.Duration) {
	c := NewParallelController(max, block, timeout)
	var oldPP = (*unsafe.Pointer)(unsafe.Pointer(&m.c))
	if p := atomic.SwapPointer(oldPP, unsafe.Pointer(c)); p != nil {
		oldC := (*ParallelController)(p)
		oldC.close()
	}
}

// Acquire
//
//	@Description: 获取并发权限
//	@receiver m
//	@return releaser 返回释放器
//	@return state 状态
func (m *ParallelLimiter) Acquire() (releaser ParallelReleaser, state PARALLEL_ACQUIRE_STATE) {
	if p := m.c; p == nil {
		panic("ParallelLimiter Acquire nil")
	} else {
		return p, p.Acquire()
	}
}

type ParallelReleaser interface {
	Release() bool
}

func NewParallelController(max int32, block bool, timeout time.Duration) *ParallelController {
	return &ParallelController{ch: make(chan struct{}, max), block: block, timeout: timeout}
}

// ParallelController 并发控制器
type ParallelController struct {
	ch      chan struct{}
	block   bool
	timeout time.Duration
}

func (m *ParallelController) close() {
	close(m.ch)
}

type PARALLEL_ACQUIRE_STATE int8

const PARALLEL_ACQUIRE_SUCCESS PARALLEL_ACQUIRE_STATE = 0 //正常
const PARALLEL_ACQUIRE_MAX PARALLEL_ACQUIRE_STATE = 1     //达到最大并发
const PARALLEL_ACQUIRE_TIMEOUT PARALLEL_ACQUIRE_STATE = 2 //超时
const PARALLEL_ACQUIRE_CLOSED PARALLEL_ACQUIRE_STATE = 3  //关闭

// Acquire
//
//	@Description: 获取
//	@receiver m
//	@return state 状态
func (m *ParallelController) Acquire() (state PARALLEL_ACQUIRE_STATE) {
	defer ExceptionError(func(e error) {
		state = PARALLEL_ACQUIRE_CLOSED
	})
	if m.block {
		if m.timeout > 0 {
			timer := time.NewTimer(m.timeout)
			select {
			case m.ch <- struct{}{}:
				return PARALLEL_ACQUIRE_SUCCESS
			case <-timer.C:
				return PARALLEL_ACQUIRE_TIMEOUT
			}
		} else {
			m.ch <- struct{}{}
			return PARALLEL_ACQUIRE_SUCCESS
		}
	} else {
		select {
		case m.ch <- struct{}{}:
			return PARALLEL_ACQUIRE_SUCCESS
		default:
			return PARALLEL_ACQUIRE_MAX
		}
	}
}

// Release
//
//	@Description: 释放
//	@receiver m
//	@return bool
func (m *ParallelController) Release() bool {
	select {
	case <-m.ch:
		return true
	default:
		return false
	}
}

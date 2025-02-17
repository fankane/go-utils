package function

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/fankane/go-utils/goroutine"
	"go.uber.org/ratelimit"
)

/**
基于漏桶算法，执行函数, 函数按 rate 频率

漏桶算法：go.uber.org/ratelimit
*/

var ErrOutOfLimit = errors.New("capacity exceeding the limit")

type DoLeaky struct {
	rate     int64 //每秒执行数量
	capacity int64

	limiter  ratelimit.Limiter
	count    atomic.Int64
	funcChan chan func()

	async bool
}

type LeakyParam struct {
	ExecASync bool
}
type LeakyOption func(param *LeakyParam)

// ASync 异步执行已经可以执行的函数
func ASync() LeakyOption {
	return func(param *LeakyParam) {
		param.ExecASync = true
	}
}

func NewLeakyFunc(rate, capacity int64, opts ...LeakyOption) (*DoLeaky, error) {
	if rate <= 0 || capacity <= 0 {
		return nil, fmt.Errorf("rate or capacity is invalid")
	}
	dp := &LeakyParam{}
	for _, opt := range opts {
		opt(dp)
	}
	d := &DoLeaky{
		rate:     rate,
		capacity: capacity,
		limiter:  ratelimit.New(int(rate)),
		count:    atomic.Int64{},
		funcChan: make(chan func(), capacity+1),
		async:    dp.ExecASync,
	}
	go d.exec()
	return d, nil
}

func (l *DoLeaky) exec() {
	defer goroutine.Recover()
	for {
		select {
		case f := <-l.funcChan:
			if l.async {
				go func() {
					defer goroutine.Recover()
					f()
				}()
			} else {
				f()
			}
			l.count.Add(-1)
		}
	}
}

// Push add function to pool, block when channel full
func (l *DoLeaky) Push(f func()) error {
	if l == nil || l.limiter == nil {
		return fmt.Errorf("leaky is nil")
	}
	if l.count.Load() > l.capacity {
		return ErrOutOfLimit
	}
	l.count.Add(1)
	l.limiter.Take()
	l.funcChan <- f
	return nil
}

// PushBlock add function to channel, block when channel full
func (l *DoLeaky) PushBlock(f func()) error {
	if l == nil || l.limiter == nil {
		return fmt.Errorf("leaky is nil")
	}
	l.count.Add(1)
	l.limiter.Take()
	l.funcChan <- f
	return nil
}

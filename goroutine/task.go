package goroutine

import (
	"context"
	"fmt"
	"github.com/fankane/go-utils/utime"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

const (
	defaultMaxWait = 100000 //默认最多等待任务：10万个
)

type TaskManager struct {
	funcChan   chan func()
	started    bool
	released   bool
	lock       sync.Mutex
	runnerPool *ants.Pool
	taskOptions
}

type taskOptions struct {
	runnerNum       int           //执行者数量
	delayStart      bool          //延迟启动
	maxWaitTask     int           //最多等待数量
	gracefulTimeout time.Duration //优雅释放时的超时时间
}

type TOption func(*TaskManager)

var defaultOpt = taskOptions{
	runnerNum:   1,
	delayStart:  false,
	maxWaitTask: defaultMaxWait,
}

func NewTaskManager(opts ...TOption) (*TaskManager, error) {
	tm := &TaskManager{
		taskOptions: defaultOpt,
	}
	for _, opt := range opts {
		opt(tm)
	}
	tm.funcChan = make(chan func(), tm.maxWaitTask)
	pool, err := ants.NewPool(tm.runnerNum)
	if err != nil {
		return nil, err
	}
	tm.runnerPool = pool

	if !tm.delayStart {
		if err = tm.Start(); err != nil {
			return nil, err
		}
	}
	return tm, nil
}

func (t *TaskManager) Start() error {
	if t.started {
		return fmt.Errorf("already started")
	}
	if t.released {
		return fmt.Errorf("has been relaesed")
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	go func() {
		for {
			if t.released {
				return
			}
			select {
			case task := <-t.funcChan:
				if t.runnerPool.IsClosed() { //已经释放，直接返回
					return
				}
				t.runnerPool.Submit(task)
			}
		}
	}()
	t.started = true
	return nil
}

func (t *TaskManager) Release() {
	if t.released {
		return
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.released { //获取到锁以后，再判断一次, 防止已经被释放过导致的异常
		return
	}
	t.released = true
	t.emptyChan()
	t.closeResource()
}

// GracefulRelease 优雅释放，不再接收新任务，同时等待未完成的任务继续完成
func (t *TaskManager) GracefulRelease() {
	if t.released {
		return
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.released { //获取到锁以后，再判断一次, 防止已经被释放过导致的异常
		return
	}

	t.released = true
	ctx := context.Background()
	if t.gracefulTimeout.Nanoseconds() > 0 {
		ctx, _ = context.WithTimeout(ctx, t.gracefulTimeout)
	}

	finish := make(chan bool)
	go func() {
		tk := time.NewTicker(time.Millisecond)
		defer tk.Stop()
		for range tk.C {
			if t.RunningCnt() == 0 {
				finish <- true
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		t.emptyChan()
		break
	case <-finish:
		break
	}
	t.closeResource()
}
func (t *TaskManager) closeResource() {
	t.runnerPool.Release()
	go func() {
		utime.DelayDo(time.Second, func() error { //延迟一秒关闭channel，防止高并发的时候，addTask 异常
			close(t.funcChan)
			return nil
		})
	}()
}

// RunningCnt 正在运行中的任务，包括执行中的和等待中的
func (t *TaskManager) RunningCnt() int {
	return t.runnerPool.Running() + len(t.funcChan)
}

func (t *TaskManager) AddTask(tasks ...func()) error {
	for _, task := range tasks {
		if t.released {
			return fmt.Errorf("task has been released")
		}
		select {
		case t.funcChan <- task:
		default:
			return fmt.Errorf("task chan is full")
		}
	}
	return nil
}

func (t *TaskManager) emptyChan() {
	if !t.released || !t.started {
		return //没释放，或者未开始的任务不清空
	}
	go func() {
		for {
			select {
			case <-t.funcChan: //清空chan 里面的消息
			}
		}
	}()
}

func WithMaxWait(max int) TOption {
	return func(manager *TaskManager) {
		manager.maxWaitTask = max
	}
}
func WithRunnerNum(runnerNum int) TOption {
	return func(manager *TaskManager) {
		manager.runnerNum = runnerNum
	}
}
func WithDelayStart(delay bool) TOption {
	return func(manager *TaskManager) {
		manager.delayStart = delay
	}
}
func WithGraceTimeout(timeout time.Duration) TOption {
	return func(manager *TaskManager) {
		manager.gracefulTimeout = timeout
	}
}

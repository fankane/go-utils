package function

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/fankane/go-utils/goroutine"
	"github.com/fankane/go-utils/plugin/log"
)

var ErrTimeout = errors.New("timeout")

func DoWithTimeout(f func() error, timeout time.Duration) error {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	var err error
	finish := make(chan bool)
	go func() {
		defer goroutine.Recover()
		err = f()
		finish <- true
	}()

	select {
	case <-ctx.Done():
		return ErrTimeout
	case <-finish:
		return err
	}
}

func DoPrintCost(ctx context.Context, f func() error, opts ...PrintOption) error {
	start := time.Now()
	defer func() {
		params := &PrintParam{}
		for _, opt := range opts {
			opt(params)
		}
		if params.traceLog != nil {
			if params.traceCtx != nil {
				params.traceLog.DebugfCtx(params.traceCtx, "%s cost:%s", params.name, time.Since(start))
			} else {
				params.traceLog.DebugfCtx(ctx, "%s cost:%s", params.name, time.Since(start))
			}
		} else {
			fmt.Println(params.name, "cost:%s", time.Since(start))
		}
	}()
	return f()
}

func DoWithLock(f func(), locker sync.Locker, opts ...LockOption) {
	params := &LockParam{locker: locker}
	for _, opt := range opts {
		opt(params)
	}
	params.locker.Lock()
	defer params.locker.Unlock()
	f()
}

type PrintParam struct {
	name     string
	traceLog *log.Log        //没有log默认使用fmt.Print打印
	traceCtx context.Context //如有值，优先使用
}

type PrintOption func(param *PrintParam)

func CostName(name string) PrintOption {
	return func(param *PrintParam) {
		param.name = name
	}
}
func TraceLog(log *log.Log) PrintOption {
	return func(param *PrintParam) {
		param.traceLog = log
	}
}
func TraceCTX(traceCtx context.Context) PrintOption {
	return func(param *PrintParam) {
		param.traceCtx = traceCtx
	}
}

type LockParam struct {
	locker sync.Locker
}
type LockOption func(param *LockParam)

func Locker(locker sync.Locker) LockOption {
	return func(param *LockParam) {
		param.locker = locker
	}
}

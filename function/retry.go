package function

import (
	"context"
	"time"

	"github.com/fankane/go-utils/goroutine"
)

type TryParam struct {
	Max      int
	Duration time.Duration
	Timeout  time.Duration
}

type Option func(param *TryParam)

func WithMax(max int) Option {
	return func(param *TryParam) {
		param.Max = max
	}
}
func WithDuration(duration time.Duration) Option {
	return func(param *TryParam) {
		param.Duration = duration
	}
}
func WithTimeout(timeout time.Duration) Option {
	return func(param *TryParam) {
		param.Timeout = timeout
	}
}

func Retry(f func() error, opts ...Option) error {
	tp := &TryParam{Max: 1}
	for _, opt := range opts {
		opt(tp)
	}

	ctx := context.Background()
	if tp.Timeout.Nanoseconds() > 0 {
		ctx, _ = context.WithTimeout(ctx, tp.Timeout)
	}

	var err error
	finish := make(chan bool)
	go func() {
		defer goroutine.Recover()
		num := 0
	RETRY:
		err = f()
		if err != nil {
			num++
			if num < tp.Max {
				time.Sleep(tp.Duration)
				goto RETRY
			}
		}
		finish <- true
	}()

	select {
	case <-ctx.Done():
		return errTimeout
	case <-finish:
		return err
	}
}

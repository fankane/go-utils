package goroutine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

const defaultMax = -1 //默认没有限制

type OptParams struct {
	Max             int           //协程最大数量, 默认无穷大
	ReturnWhenError bool          //当出错的时候返回, 默认不返回
	Timeout         time.Duration //超时后结束，默认没有超时时间
	DisableRecovery bool          //不捕获panic
}

type Option func(params *OptParams)

/*
*
Exec 并发执行 function
可指定最大协程数量，超时时间，错误是否返回
*/
func Exec(fs []func() error, opts ...Option) error {
	if len(fs) == 0 {
		return fmt.Errorf("empty func to exec")
	}
	params := &OptParams{Max: defaultMax}
	for _, opt := range opts {
		opt(params)
	}
	pool, err := ants.NewPool(params.Max)
	if err != nil {
		return err
	}
	defer pool.Release()

	errChan := make(chan error)
	finish := make(chan bool)
	go func() {
		if !params.DisableRecovery {
			defer Recover()
		}
		wg := sync.WaitGroup{}
		for i := 0; i < len(fs); i++ {
			wg.Add(1)
			idx := i
			if er := pool.Submit(func() {
				defer wg.Done()
				tempErr := fs[idx]()
				if params.ReturnWhenError && tempErr != nil {
					errChan <- tempErr
				}
			}); er != nil {
				errChan <- er
			}
		}
		wg.Wait()
		finish <- true
	}()

	ctx := context.Background()
	if params.Timeout.Nanoseconds() > 0 {
		ctx, _ = context.WithTimeout(ctx, params.Timeout)
	}
	select {
	case e := <-errChan:
		return e
	case <-ctx.Done():
		return ctx.Err()
	case <-finish:
		return nil
	}
}

func WithMax(max int) Option {
	return func(params *OptParams) {
		params.Max = max
	}
}
func WithReturnWhenError(returnWhenError bool) Option {
	return func(params *OptParams) {
		params.ReturnWhenError = returnWhenError
	}
}
func WithTimeout(timeout time.Duration) Option {
	return func(params *OptParams) {
		params.Timeout = timeout
	}
}

func DisableRecovery(disable bool) Option {
	return func(params *OptParams) {
		params.DisableRecovery = disable
	}
}

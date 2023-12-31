package function

import (
	"context"
	"errors"
	"time"

	"github.com/fankane/go-utils/goroutine"
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

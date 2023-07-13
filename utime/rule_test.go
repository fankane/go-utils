package utime

import (
	"fmt"
	"testing"
	"time"
)

func TestTickerDo(t *testing.T) {
	err := TickerDo(time.Millisecond*500, func() error {
		fmt.Println(time.Now())
		time.Sleep(time.Second)
		return nil
	}, WithMax(10), WithReturn(true), WithDoExactly(true), WithFirstImmediately(true))
	if err != nil {
		fmt.Println("err", err)
		return
	}
}

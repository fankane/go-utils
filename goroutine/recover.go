package goroutine

import (
	"fmt"
	"runtime"
)

func Recover() {
	if err := recover(); err != nil {
		var buf [2048]byte
		n := runtime.Stack(buf[:], false)
		fmt.Printf("[panic] err: %v\nstack: %s\n", err, buf[:n])
	}
}

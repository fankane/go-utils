package limiter

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewSlidingWindowLimiter(t *testing.T) {
	l := NewSlidingWindowLimiter(time.Millisecond, 100)

	sg := &sync.WaitGroup{}
	for a := 0; a < 3; a++ {
		for i := 0; i < 20; i++ {
			if i == 10 {
				time.Sleep(time.Millisecond * 500)
				if a == 1 {
					time.Sleep(time.Second)
				}
			}

			sg.Add(1)
			t1 := i
			go func(a int) {
				defer sg.Done()
				r := l.Allow()
				if !r {
					fmt.Println("out of limiter,", a)
				}
			}(t1)
		}
		time.Sleep(time.Second)
		fmt.Println("----------------------")
	}

	sg.Wait()
}

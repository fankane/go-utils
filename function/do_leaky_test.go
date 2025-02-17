package function

import (
	"fmt"
	"testing"
	"time"
)

func TestDoLeaky_Push(t *testing.T) {
	l, err := NewLeakyFunc(50, 20)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	for i := 0; i < 50; i++ {
		t1 := i
		go func(t1 int) {
			if err := l.PushBlock(func() {
				fmt.Println("i am:", t1, ", time:", time.Now())
			}); err != nil {
				fmt.Println("err:", err)
			}
		}(t1)
		//l.Push(func() {
		//	fmt.Println("i am:", t1, ", time:", time.Now())
		//})
	}
	fmt.Println("tt:", time.Now())
	time.Sleep(time.Minute)
}

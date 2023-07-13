package goroutine

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTaskManager(t *testing.T) {
	tm, err := NewTaskManager(WithGraceTimeout(time.Second), WithRunnerNum(10))
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	go func() {
		for i := 0; i < 50; i++ {
			temp := i
			//fmt.Println("add task i=", i)
			if err2 := tm.AddTask(func() {
				fmt.Println("i=", temp)
				time.Sleep(time.Millisecond * time.Duration(i+1) * 10)
			}); err2 != nil {
				fmt.Println("add task err:", err2)
				return
			}
			//time.Sleep(time.Millisecond * 100)
		}
	}()
	//time.Sleep(time.Millisecond * 50)
	go func() {
		tk := time.NewTicker(time.Millisecond * 3)
		for range tk.C {
			fmt.Println("running:", tm.RunningCnt())
			if tm.RunningCnt() == 0 {
				return
			}
		}
	}()

	time.Sleep(time.Millisecond * 100)

	go func() {
		fmt.Println("GracefulRelease...")
		//tm.GracefulRelease()
		tm.Release()
	}()
	go func() {
		for i := 0; i < 10; i++ {
			if err := tm.AddTask(func() {
				fmt.Println("hello ...")
			}); err != nil {
				fmt.Println("err:", err)
			}
			time.Sleep(time.Millisecond)
		}
	}()

	time.Sleep(time.Second * 5)
}

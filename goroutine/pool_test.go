package goroutine

import (
	"fmt"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	fs := make([]func() error, 0)
	for i := 0; i < 10; i++ {
		temp := i
		fs = append(fs, func() error {
			if temp > 5 {
				time.Sleep(time.Second)
				return fmt.Errorf("test")
			}
			time.Sleep(time.Second)
			fmt.Println("i=", temp)
			return nil
		})
	}
	start := time.Now()
	for i := 0; i < 10; i++ {
		err := Exec(fs, WithReturnWhenError(true), WithMax(50))
		if err != nil {
			fmt.Println("err:", err)
		}
		fmt.Println("cost:", time.Since(start))
		fmt.Println("----------------------")
	}

}

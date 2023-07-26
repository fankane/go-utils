package utime

import (
	"fmt"
	"testing"
	"time"
)

func TestCronDo(t *testing.T) {
	if err := CronDo("*/1 * * * *", func() {
		//if err := CronDo("*/5 * * * *", func() {
		fmt.Println(time.Now())
	}); err != nil {
		fmt.Println("err:", err)
		return
	}
	time.Sleep(time.Minute * 5)
}

package function

import (
	"fmt"
	"testing"
	"time"

	"github.com/fankane/go-utils/utime"
)

func TestRetry(t *testing.T) {
	f := func() error {
		fmt.Println(time.Now().Format(utime.LayYMDHms1))
		time.Sleep(time.Second)
		return fmt.Errorf("test")
	}

	err := Retry(f, WithMax(5), WithDuration(time.Second), WithTimeout(time.Second*4))
	if err != nil {
		fmt.Println("err:", err)
		return
	}
}

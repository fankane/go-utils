package machine

import (
	"fmt"
	"testing"
	"time"

	"github.com/fankane/go-utils/str"
)

func TestGetCPUPercent(t *testing.T) {
	fmt.Println(GetCPUPercent())
	fmt.Println(GetMemPercent())
	fmt.Println(str.FormatFileSize(float64(GetMemUsed())))
	fmt.Println(str.FormatFileSize(float64(GetMemAvailable())))
	fmt.Println(str.FormatFileSize(float64(GetMemTotal())))
	fmt.Println(GetDiskUsedPercent())
}

func TestGetSelfCPUPercent(t *testing.T) {
	go func() {
		res := 1
		for i := 0; i < 100; i++ {
			res += i
			time.Sleep(time.Millisecond * 300)
		}
		fmt.Println(res)
	}()
	time.Sleep(time.Second)
	for i := 0; i < 10; i++ {
		fmt.Println(GetSelfCPUPercent())
		fmt.Println(str.FormatFileSize(float64(GetSelfMemory())))
		fmt.Println("-------------")
		time.Sleep(time.Second)
	}

}

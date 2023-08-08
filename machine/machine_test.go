package machine

import (
	"fmt"
	"github.com/fankane/go-utils/str"
	"testing"
)

func TestGetCPUPercent(t *testing.T) {
	fmt.Println(GetCPUPercent())
	fmt.Println(GetMemPercent())
	fmt.Println(str.FormatFileSize(float64(GetMemUsed())))
	fmt.Println(str.FormatFileSize(float64(GetMemAvailable())))
	fmt.Println(str.FormatFileSize(float64(GetMemTotal())))
	fmt.Println(GetDiskUsedPercent())
}

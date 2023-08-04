package pprof

import (
	"fmt"
	"testing"
	"time"

	"github.com/fankane/go-utils/plugin"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	time.Sleep(time.Minute)
}

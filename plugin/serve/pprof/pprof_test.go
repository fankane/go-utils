package pprof

import (
	"fmt"
	"github.com/fankane/go-utils/plugin"
	"testing"
	"time"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	time.Sleep(time.Minute)
}

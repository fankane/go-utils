package log

import (
	"fmt"
	"github.com/fankane/go-utils/plugin"
	"testing"
)

func Test_newLogger(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	if Logger == nil {
		fmt.Println("err")
		return
	}
	for i := 0; i < 9; i++ {
		Logger.Info("test", i)
		GetLogger("logName2").Errorf("hi")
	}
}

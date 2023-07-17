package log

import (
	"fmt"
	"testing"
)

func Test_newLogger(t *testing.T) {
	if Logger == nil {
		fmt.Println("err")
		return
	}
	for i := 0; i < 100; i++ {
		Logger.Info("test", i)
	}

}

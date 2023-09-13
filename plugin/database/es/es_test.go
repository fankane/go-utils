package es

import (
	"fmt"
	"github.com/fankane/go-utils/plugin"
	"testing"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	if Cli == nil {
		fmt.Println("es client is nil")
		return
	}
	fmt.Println("success")
}

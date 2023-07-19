package mysql

import (
	"fmt"
	"testing"

	"github.com/fankane/go-utils/plugin"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	if DB == nil {
		fmt.Println("db is nil")
		return
	}
}

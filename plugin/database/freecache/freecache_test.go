package freecache

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
	if Cache == nil {
		fmt.Println("Cache is nil")
		return
	}
	Cache.Set([]byte("test"), []byte("hello"), 100)
	res, err := Cache.Get([]byte("test"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("res:", string(res))
}

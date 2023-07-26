package conf

import (
	"fmt"
	"testing"
	"time"

	"github.com/fankane/go-utils/plugin"
	"github.com/fankane/go-utils/str"
	"github.com/fankane/go-utils/utime"
)

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}

	x := &AB{}
	if err := Unmarshal(x); err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(str.ToJSON(x))

	utime.TickerDo(time.Second*3, func() error {
		fmt.Println(time.Now(), str.ToJSON(x))
		return nil
	})

	time.Sleep(time.Minute)
}

type AB struct {
	A int    `yaml:"a"`
	B string `yaml:"b"`
}

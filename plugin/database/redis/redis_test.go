package redis

import (
	"context"
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
	if Client == nil {
		fmt.Println("redis client is nil")
		return
	}
	Client.Set(context.Background(), "hf", time.Now().String(), 0)
	res, err := Client.Get(context.Background(), "hf").Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

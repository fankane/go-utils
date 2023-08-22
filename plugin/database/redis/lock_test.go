package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/fankane/go-utils/plugin"
)

func Test_rdsLock_Lock(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	if Client == nil {
		fmt.Println("redis client is nil")
		return
	}
	//rl := NewRdsLock(Client)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			rl := NewRdsLock(Client)
			ok, err := rl.Lock("test")
			if err != nil {
				fmt.Println("lock err:", err)
				return
			}
			if ok {
				fmt.Println("get lock idx:", idx)
			} else {
				fmt.Println("cannot lock idx:", idx)
			}
			o2, err := rl.Release()
			if err != nil {
				fmt.Println("lock err:", err)
				return
			}
			if o2 {
				fmt.Println("released idx:", idx)
			} else {
				fmt.Println("cannot released idx:", idx)
			}
		}(i)
	}
	time.Sleep(time.Second * 3)
}

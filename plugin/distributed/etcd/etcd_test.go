package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fankane/go-utils/plugin"
)

type testH struct {
}

/**
PutHandle(key, value []byte, version int64)
DelHandle(key []byte)
*/
func (h testH) PutHandle(key, value []byte, version int64) {
	fmt.Println(fmt.Sprintf("put handle>>key:%s, val:%s, version:%d", string(key), string(value), version))
}
func (h testH) DelHandle(key []byte) {
	fmt.Println(fmt.Sprintf("del handle>>key:%s", string(key)))
}

func TestFactory_Setup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}

	if Op == nil {
		fmt.Println("op is nil")
		return
	}
	ctx := context.Background()

	go func() {
		for i := 0; i < 10; i++ {
			_, err := Op.Put(ctx, "k1", time.Now().String())
			if err != nil {
				fmt.Println("put err:", err)
				return
			}
			time.Sleep(time.Millisecond * 500)
		}
		time.Sleep(time.Second * 10)
		for i := 0; i < 10; i++ {
			_, err := Op.Put(ctx, "k1", time.Now().String())
			if err != nil {
				fmt.Println("put err:", err)
				return
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
	go func() {
		time.Sleep(time.Second * 7)
		deled, err := Op.Delete(ctx, "k1")
		if err != nil {
			fmt.Println("del err:", err)
			return
		}
		fmt.Println("deleted:", deled)
	}()

	resM, err := Op.Get(ctx, "k1")
	if err != nil {
		fmt.Println("get err:", err)
		return
	}

	for s, v := range resM {
		fmt.Println(s, ", value:", string(v.Val), v.Version)
	}

	fmt.Println("--------------------------------------")
	fmt.Println("--------------------------------------")
	Op.Watch(ctx, "k1", &testH{})

	time.Sleep(time.Minute * 5)
}

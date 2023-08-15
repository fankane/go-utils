package nsq

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
	if DefaultProducer == nil {
		fmt.Println("producer is nil")
		return
	}
	go func() {
		for i := 0; i < 30; i++ {
			if err := DefaultProducer.SendMsg("test1", []byte(fmt.Sprintf("hello:%d", i+1))); err != nil {
				fmt.Println("publish err:", err)
				continue
			}
			fmt.Println("publish success")
			time.Sleep(time.Second * 3)
		}
	}()
	go func() {
		for i := 0; i < 30; i++ {
			if err := GetProducer("p2").SendMsg("test1", []byte(fmt.Sprintf("hi:%d", i+1))); err != nil {
				fmt.Println("publish err:", err)
				continue
			}
			fmt.Println("publish success")
			time.Sleep(time.Second * 3)
		}
	}()

	if err := RegisterHandler("c1", func(ctx context.Context, value []byte) error {
		fmt.Println(fmt.Sprintf("Attempts:%d, NSQDAddress:%s, time:%d", Attempts(ctx),
			NSQDAddress(ctx), Timestamp(ctx)), "value:", string(value))
		return nil
	}); err != nil {
		fmt.Println("consumer err:", err)
		return
	}
	if err := RegisterHandler("c2", func(ctx context.Context, value []byte) error {
		fmt.Println(fmt.Sprintf("c2 Attempts:%d, NSQDAddress:%s, time:%d", Attempts(ctx),
			NSQDAddress(ctx), Timestamp(ctx)), "value:", string(value))
		return nil
	}); err != nil {
		fmt.Println("consumer err:", err)
		return
	}

	time.Sleep(time.Minute)
}

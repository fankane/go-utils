package rabbit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fankane/go-utils/plugin"
)

func TestNewProducer(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	if DefaultProducer == nil {
		fmt.Println("producer is nil")
		return
	}
	go func() {
		time.Sleep(time.Second * 15)
		for i := 0; i < 10; i++ {
			if err := DefaultProducer.SendMsg(context.Background(), "test1", []byte(fmt.Sprintf("aaaa - %d", i+1))); err != nil {
				fmt.Println("err:", err)
				return
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
	go func() {
		time.Sleep(time.Second * 25)
		for i := 0; i < 10; i++ {
			if err := DefaultProducer.SendMsg(context.Background(), "test2", []byte(fmt.Sprintf("bbbb - %d", i+1))); err != nil {
				fmt.Println("err:", err)
				return
			}
			time.Sleep(time.Millisecond * 600)
		}
	}()

	go func() {
		if err := RegisterHandler("c1", func(ctx context.Context, value []byte) error {
			fmt.Println("receive:", string(value))
			return nil
		}); err != nil {
			fmt.Println("err:", err)
			return
		}
	}()
	time.Sleep(time.Minute * 3)
	fmt.Println("success")
}

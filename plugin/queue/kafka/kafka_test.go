package kafka

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fankane/go-utils/plugin"
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
	fmt.Println("producer init success")
	DefaultProducer.SendMessage("test1", []byte(fmt.Sprintf("key-%d", time.Now().Unix())), []byte("hello world !!!"))
	time.Sleep(time.Second * 2)
	if err := RegisterHandler("c1", func(ctx context.Context, key, value []byte) error {
		fmt.Println(fmt.Sprintf("t:%s, p:%d, o:%d, ts:%s", Topic(ctx), Partition(ctx), Offset(ctx), Timestamp(ctx)),
			"business key:", string(key), "value:", string(value))
		return nil
	}); err != nil {
		fmt.Println("consumer err:", err)
		return
	}
	fmt.Println("register handler success")
	go func() {
		time.Sleep(time.Second * 5)
		for i := 0; i < 3; i++ {
			DefaultProducer.SendMessage("test1", []byte(fmt.Sprintf("key-%d", time.Now().Unix())), []byte(fmt.Sprintf("val:%d", i+100)))
			time.Sleep(time.Second * 2)
		}
	}()

	time.Sleep(time.Second * 66)
}

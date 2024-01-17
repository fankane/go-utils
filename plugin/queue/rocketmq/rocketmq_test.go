package rocketmq

import (
	"context"
	"fmt"
	"github.com/fankane/go-utils/str"
	"testing"
	"time"

	"github.com/fankane/go-utils/plugin"
)

func loadPlugin() error {
	if DefaultProducer == nil {
		return plugin.Load()
	}
	return nil
}

func TestFactory_Producer(t *testing.T) {
	loadPlugin()
	for i := 30; i < 40; i++ {
		res, err := DefaultProducer.SendSync(context.Background(), "test_topic_hf3",
			[]byte(fmt.Sprintf("test2 for producer of %d", i+1)), WithTag("qqq"))
		if err != nil {
			fmt.Println("SendSync err:", err)
			continue
		}
		fmt.Println(str.ToJSON(res))
		time.Sleep(time.Second * 5)
	}

}
func TestFactory_Consumer(t *testing.T) {
	loadPlugin()
	if DefaultProducer == nil {
		fmt.Println("producer is nil")
		return
	}
	if err := RegisterHandler("c1", func(ctx context.Context, value []byte) error {
		ts := GetStoreTimestamp(ctx)
		fmt.Println("ts format:", time.UnixMilli(ts), "content:", string(value))
		//fmt.Println("consume:", string(value))
		return nil
	}); err != nil {
		fmt.Println("RegisterHandler err:", err)
		return
	}
	fmt.Println("RegisterHandler success")
	time.Sleep(time.Hour)
}

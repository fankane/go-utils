package rocketmq

import (
	"context"
	"fmt"
	"testing"
)

func TestConsumer_Start(t *testing.T) {
	_, err := NewConsumer(&ConsumerConf{
		NameServerAddrs: []string{"192.168.99.38:9876"},
		GroupName:       "hello_fan",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = DeleteTopic(context.Background(), []string{"192.168.99.38:9876"}, "192.168.99.38:10911", "auto_create2")
	//err = c.Start(AutoCreateTopic(true), TopicHandler("auto_create10", "192.168.99.38:10911", func(ctx context.Context, value []byte) error {
	//	fmt.Println("value:", string(value))
	//	return nil
	//}))
	if err != nil {
		fmt.Println("start err:", err)
		return
	}
}

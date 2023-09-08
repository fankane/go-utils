package memory

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
	topic := "hufan"
	//go func() {
	//	//time.Sleep(time.Second * 2)
	//	for i := 0; i < 10; i++ {
	//		if err := NewProducer().SendMessage(topic, []byte(fmt.Sprintf("i=%d", i))); err != nil {
	//			fmt.Println("err:", err, ", i=", i)
	//			return
	//		}
	//		//time.Sleep(time.Second)
	//	}
	//}()
	//go func() {
	//	time.Sleep(time.Second * 2)
	//	for i := 0; i < 10; i++ {
	//		if err := NewProducer().SendMessage(topic, []byte(fmt.Sprintf("i=%d", i))); err != nil {
	//			fmt.Println("err:", err)
	//			//return
	//		}
	//		time.Sleep(time.Millisecond * 300)
	//	}
	//}()
	//
	//time.Sleep(time.Second * 3)
	go func() {
		//time.Sleep(time.Second * 13)
		for i := 0; i < 10; i++ {
			if i < 4 {
				if err := NewProducer().SendMessage(topic+"_copy", []byte(fmt.Sprintf("i=%d", i+100)), Delay(time.Second*time.Duration(i+1))); err != nil {
					fmt.Println("err:", err)
					return
				}
			} else if i == 8 || i == 7 {
				if err := NewProducer().SendMessage(topic+"_copy", []byte(fmt.Sprintf("i=%d", i+100)), Delay(time.Second*2)); err != nil {
					fmt.Println("err:", err)
					return
				}
			} else {
				if err := NewProducer().SendMessage(topic+"_copy", []byte(fmt.Sprintf("i=%d", i+100))); err != nil {
					fmt.Println("err:", err)
					return
				}
			}
			//time.Sleep(time.Second)
		}
	}()

	RegisterHandler(topic, func(ctx context.Context, value []byte) error {
		fmt.Println("消费数据:", time.Now(), string(value))
		return nil
	})
	//time.Sleep(time.Second)
	RegisterHandler(topic+"_copy", func(ctx context.Context, value []byte) error {
		fmt.Println("copy 消费数据:", time.Now(), string(value))
		return nil
	})
	time.Sleep(time.Minute)
}

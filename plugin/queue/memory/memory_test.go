package memory

import (
	"context"
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

func TestStopConsumer(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	topic := "test1"
	go func() {
		utime.TickerDo(time.Millisecond*1000, func() error {
			fmt.Println("left:", CachedLen(), ", cacheSize:", str.FormatFileSize(float64(CachedSize())))
			return nil
		})
	}()
	go func() {
		for i := 0; i < 10; i++ {
			if err := NewProducer().SendMessage(topic, []byte(fmt.Sprintf("i=%d", i+1))); err != nil {
				fmt.Println("err:", err)
				return
			}
			time.Sleep(time.Second)
		}
		fmt.Println("total len:", CachedLen())
	}()

	go func() {
		time.Sleep(time.Second * 3) //模拟3秒后停止消费
		if err := StopConsumer(topic); err != nil {
			fmt.Println("stop err:", err)
			return
		}
		fmt.Println("stop consumer success")
	}()
	//time.Sleep(time.Second)
	RegisterHandler(topic, func(ctx context.Context, value []byte) error {
		fmt.Println("消费数据:", time.Now(), string(value))
		return nil
	})

	go func() {
		time.Sleep(time.Second * 6) // 6秒后再次消费
		RegisterHandler(topic, func(ctx context.Context, value []byte) error {
			fmt.Println("second 消费数据:", time.Now(), string(value))
			return nil
		})
	}()
	go func() {
		if er := RegisterHandler(topic, func(ctx context.Context, value []byte) error {
			fmt.Println("second 消费数据:", time.Now(), string(value))
			return nil
		}); er != nil {
			fmt.Println("注册消费者失败 err:", er)
		}
	}()

	time.Sleep(time.Minute)
}

func TestBackup(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	topic := "test1"
	go func() {
		for i := 0; i < 15; i++ {
			if err := NewProducer().SendMessage(topic, []byte(fmt.Sprintf("i=%d", i+1))); err != nil {
				fmt.Println("err:", err)
				return
			}
		}
		fmt.Println("total len:", CachedLen())
	}()
	time.Sleep(time.Second)
	if err := StopAndBackOnFile("./backup"); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("back success")
}

func TestLoadFileData(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}
	RegisterHandler("test1", func(ctx context.Context, value []byte) error {
		pushTS := ctx.Value(CtxPushTs)
		fmt.Println("消费数据:", time.Now(), string(value), time.Unix(0, pushTS.(int64)))
		return nil
	})
	utime.TickerDo(time.Second, func() error {
		fmt.Println("left:", CachedLen())
		return nil
	})
	time.Sleep(time.Minute)
}

package etcd

import (
	"fmt"
	"testing"
	"time"

	"github.com/fankane/go-utils/plugin"
	"github.com/fankane/go-utils/str"
	"github.com/fankane/go-utils/utime"
)

func Test_etcd_RegisterServer(t *testing.T) {
	if err := plugin.Load(); err != nil {
		fmt.Println("err:", err)
		return
	}

	if Op == nil {
		fmt.Println("op is nil")
		return
	}
	go func() {
		utime.TickerDo(time.Second, func() error {
			res := Op.GetServers()
			fmt.Println("get servers:", str.ToJSON(res))
			return nil
		})
	}()

	time.Sleep(time.Second * 5)
	fmt.Println("unregister ...")
	if err := Op.UnRegisterServer(); err != nil {
		fmt.Println("err:", err)
		return
	}

	//go func() {
	//	time.Sleep(time.Second * 3)
	//	if err := Op.RegisterServer(); err != nil {
	//		fmt.Println("222 err:", err)
	//		return
	//	}
	//}()
	time.Sleep(time.Minute)
}

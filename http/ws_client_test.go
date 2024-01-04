package http

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/fankane/go-utils/str"
	"github.com/fankane/go-utils/utime"
)

var (
	testHost = "192.168.0.93:9999"
	testPath = "/ws/v1"
)

func TestWSServer(t *testing.T) {
	http.HandleFunc(testPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(str.ToJSON(r.Header))
		cli, err := ServerHandleWS(HandleWSParam{
			W: w,
			R: r,
		})
		if err != nil {
			fmt.Println("handle ws err:", err)
			return
		}
		go func() {
			time.Sleep(time.Second * 3)
			if err := cli.WriteMessage(TextMessage, []byte("我是胡帆...")); err != nil {
				fmt.Println("111 err:", err)
				return
			}
			time.Sleep(time.Second * 5)
			if err := cli.WriteMessage(TextMessage, []byte("222 我是胡帆...")); err != nil {
				fmt.Println(time.Now(), "222 err:", err)
				return
			}
		}()
		fmt.Println("ServerHandleWS success")
	})
	fmt.Println("server start success...")
	if err := http.ListenAndServe(testHost, nil); err != nil {
		fmt.Println("listen err:", err)
		return
	}
	fmt.Println("server end")
}

func TestNewWSClient(t *testing.T) {
	//testAddr := "127.0.0.1"
	u := url.URL{
		Scheme: "ws",
		Path:   testPath,
		Host:   testHost,
		//Host: "192.168.0.93:9001",
		//Path: "/g_hf_management/chat/ws/spark",
	}

	h := http.Header{}
	h.Add("Authorization", "xxx")
	cliInfo, err := NewWSClient(u, DisablePingTest(true), RequestHeader(h))
	if err != nil {
		fmt.Println("NewWSClient err:", err)
		return
	}
	fmt.Println("client create success")
	go func() {
		if err := cliInfo.WriteMessage(TextMessage, []byte(fmt.Sprintf("今天天气怎么样？"))); err != nil {
			fmt.Println("111 write err:", err)
		}

		time.Sleep(time.Second * 5)
		fmt.Println("发送第二次。。。")
		if err := cliInfo.WriteMessage(TextMessage, []byte(fmt.Sprintf("武汉工程大学各学院？"))); err != nil {
			fmt.Println("222 write err:", err)
		}
	}()
	go func() {
		cliInfo.ListenMessage(func(mt int, data []byte) {
			fmt.Println(string(data))
		})
	}()
	time.Sleep(time.Hour)
	return
	for i := 0; i < 4; i++ {
		if err = cliInfo.WriteMessage(TextMessage,
			[]byte(fmt.Sprintf("client time:%s", time.Now().Format(utime.LayYMDHms1)))); err != nil {
			fmt.Println("client write msg err:", err)
			break
		}
		fmt.Println("send success")
		time.Sleep(time.Millisecond * 900)
	}
	fmt.Println(time.Now(), "client close")
	if err = cliInfo.Close(); err != nil {
		fmt.Println("write close msg err:", err)
		return
	}
	//if err = cliInfo.WriteMessage(TextMessage, []byte("close")); err != nil {
	//	fmt.Println("write close msg err:", err)
	//	return
	//}
	if err = cliInfo.WriteMessage(TextMessage, []byte("close")); err != nil {
		fmt.Println("write close msg err:", err)
		return
	}
	//fmt.Println("write close message")
	time.Sleep(time.Second * 5)
}

func serverFunc(ctx context.Context, messageType int, p []byte) (needResponse, closeConn bool, body []byte) {

	//fmt.Println("receive msg type:", messageType)
	fmt.Println("server receive msg body:", string(p))
	//time.Sleep(time.Second)
	close := false
	if string(p) == "close" {
		close = true
	}
	return true, close, []byte(fmt.Sprintf("server time:%s", time.Now().Format(utime.LayYMDHms1)))
}

func clientFunc(ctx context.Context, messageType int, p []byte) (needResponse, closeConn bool, body []byte) {
	//fmt.Println("client receive msg type:", messageType)
	//fmt.Println("client receive msg body:", string(p))
	fmt.Println(string(p))
	return false, false, nil
	//time.Sleep(time.Second)
	close := false
	if string(p) == "close" {
		close = true
	}
	return false, close, []byte(fmt.Sprintf("client time:%s", time.Now().Format(utime.LayYMDHms1)))
}

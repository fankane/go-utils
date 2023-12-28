package http

import (
	"context"
	"fmt"
	"github.com/fankane/go-utils/utime"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var (
	testHost = "192.168.0.93:9999"
	testPath = "/ws/v1"
)

func TestWSServer(t *testing.T) {
	http.HandleFunc(testPath, func(w http.ResponseWriter, r *http.Request) {
		cli, err := ServerHandleWS(HandleWSParam{
			W: w,
			R: r,
			F: serverFunc,
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
		Host:   testHost,
		Path:   testPath,
	}

	cliInfo, err := NewWSClient(u, clientFunc, DisablePingTest(true))
	if err != nil {
		fmt.Println("NewWSClient err:", err)
		return
	}
	fmt.Println("client create success")
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
	fmt.Println("client receive msg body:", string(p))
	//time.Sleep(time.Second)
	close := false
	if string(p) == "close" {
		close = true
	}
	return false, close, []byte(fmt.Sprintf("client time:%s", time.Now().Format(utime.LayYMDHms1)))
}

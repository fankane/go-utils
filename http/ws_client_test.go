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
		//Path:   testPath,
		//Host:   testHost,
		Host: "192.168.99.45:9002",
		Path: "/g_hf_management/chat/ws/spark",
	}
	query := url.Values{}
	query.Set("user_name", "hufan")
	u.RawQuery = query.Encode()
	fmt.Println(u.String())
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

		cliInfo.ListenMessage(context.Background(), func(ctx context.Context, mt int, data []byte) {
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

func TestNewWSClient2(t *testing.T) {
	//testAddr := "127.0.0.1"
	u := url.URL{
		Scheme: "ws",
		Host:   "192.168.99.45:9001",
		Path:   "/big_language_model/v1/user/open_window",
	}

	//h := http.Header{}
	//h.Add("Authorization", "xxx")
	query := url.Values{}
	query.Set("user_name", "admin")
	query.Set("session_id", "b13641ba-98ad-492e-8ae2-ccaf298354a2")
	u.RawQuery = query.Encode()
	cliInfo, err := NewWSClient(u, DisablePingTest(true), HandshakeTimeout(time.Second*5))
	if err != nil {
		fmt.Println("NewWSClient err:", err)
		return
	}
	fmt.Println("client create success")
	go func() {
		dd := `{
    "user_name":"test",
    "plugin":"base",
    "content":"diit是哪个公司缩写",
    "session_id":"b13641ba-98ad-492e-8ae2-ccaf298354a2"
}`

		if err := cliInfo.WriteMessage(TextMessage, []byte(fmt.Sprintf(dd))); err != nil {
			fmt.Println("111 write err:", err)
		}
	}()
	go func() {
		ctx := context.Background()

		//var sID string
		//idx := 0
		cliInfo.ListenMessage(ctx, func(ctx context.Context, mt int, data []byte) {
			fmt.Println(string(data))
			//s := &SSS{}
			//json.Unmarshal(data, s)
			//sID = s.SessionID
			//if s.Answer.Status == 2 {
			//	idx = 1
			//}
		})
		//		for idx <= 0 {
		//			time.Sleep(time.Second)
		//		}
		//		c2 := `{
		//    "user_name":"test",
		//    "plugin":"base",
		//    "content":"用途有哪些",
		//    "session_id":"%s"
		//}`
		//		c2 = fmt.Sprintf(c2, sID)
		//		fmt.Println("开始二次问问题：c2:", c2)
		//		if err := cliInfo.WriteMessage(TextMessage, []byte(fmt.Sprintf(c2))); err != nil {
		//			fmt.Println("111 write err:", err)
		//		}
	}()
	time.Sleep(time.Hour)
	return

}

type SSS struct {
	SessionID string `json:"session_id"`
	Answer    struct {
		Status  int    `json:"status"`
		Seq     int    `json:"seq"`
		Content string `json:"content"`
	} `json:"answer"`
}

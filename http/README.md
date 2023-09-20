### websocket 使用demo

#### server 端
```go
// main.go
func main() {
    http.HandleFunc(testPath, func(w http.ResponseWriter, r *http.Request) {
        if srvConn, err := ServerHandleWS(HandleWSParam{
            W: w,
            R: r,
            F: serverFunc,
        }); err != nil {
            return
        }
    })
    http.ListenAndServe("127.0.0.1:1234", nil)
}

// 主动发送消息 
srvConn.WriteMessage(TextMessage, []byte("hello"))

// 业务处理, 打印客户端发送的消息，返回服务端时间字符串
func serverFunc(ctx context.Context, messageType int, p []byte) (needResponse, closeConn bool, body []byte) {
    fmt.Println("server receive msg body:", string(p))
    return true, false, []byte(fmt.Sprintf("server time:%s", time.Now().Format(utime.LayYMDHms1)))
}
```

#### client 端
```go

u := url.URL{
    Scheme: "ws",
    Host:   "127.0.0.1:1234",
    Path:   "echo",
}
cliInfo, err := NewWSClient(u, clientFunc)
if err != nil {
    fmt.Println("NewWSClient err:", err)
    return
}

// 主动发送消息
cliInfo.WriteMessage(TextMessage, []byte("hello"))

// 业务处理, 打印服务端发送的消息，返回客户端端时间字符串
func clientFunc(ctx context.Context, messageType int, p []byte) (needResponse, closeConn bool, body []byte) {
    fmt.Println("server receive msg body:", string(p))
    return true, false, []byte(fmt.Sprintf("client time:%s", time.Now().Format(utime.LayYMDHms1)))
}
```


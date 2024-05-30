### Client
```go
// 发起 SSE请求，服务端返回的消息放在channel里面，
statusCode, resChan, err := NewClient().SSEGet(ctx, url)

for e := range resChan {
    if msg == nil {
        return
    }
    fmt.Println(string(e.Data))
}


```

### Server
```go
import (
    "github.com/fankane/go-utils/http"
)

// 实现SSE接口
dataChan := make(chan []byte, 100)
var err error
go func() {
    for {
		dataChan <- "hello" //需要写回Client的数据
		time.sleep(time.Second)
    }
}()

http.RegisterSSE(ResponseWriter, Request, dataChan)

```

### 部署建议
> 由于大部分场景，后台服务都有使用Nginx做代理，会发现SSE的接口没有返回数据，可配置对应location如下

```shell
location = /sse/full/path {
    # 禁用代理缓冲  
    proxy_buffering off;  
    # 如果使用SSL，还需要确保SSL会话不被缓存  
    proxy_ssl_session_reuse off;
    # 适量设置较长的超时时间
    proxy_read_timeout 86400s;

    # 原始配置 ... 
}
```

